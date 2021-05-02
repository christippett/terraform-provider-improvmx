package improvmx

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

func setupClient(t *testing.T) Client {
	apiKey := os.Getenv("IMPROVMX_API_KEY")
	if apiKey == "" {
		t.Fatal("'IMPROVMX_API_KEY' must be set for tests")
	}

	httpClient := http.Client{
		Timeout: time.Second * 20,
	}

	return NewClient(
		"https://api.improvmx.com/v3",
		apiKey,
		&httpClient,
		os.Stdout,
	)
}

func TestIntegration_Account(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	_, err := c.GetAccount(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_ListDomains(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	domain, err := c.AddDomain(ctx, &Domain{Domain: "example.com"})
	if err != nil {
		t.Fatal(err)
	}

	domains, err := c.ListDomains(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	defer c.DeleteDomain(ctx, domain)

	if len(*domains) < 1 {
		t.Errorf("domain list returned unexpected count")
	}
}

func TestIntegration_DomainCRUD(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	domain, err := c.AddDomain(ctx, &Domain{
		Domain:            "example.com",
		NotificationEmail: "test@christippett.dev",
	})
	if err != nil {
		t.Fatal(err)
	}

	domain.NotificationEmail = "test+updated@christippett.dev"
	_, err = c.UpdateDomain(ctx, domain)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := c.GetDomain(ctx, "example.com")
	if err != nil {
		t.Fatal(err)
	}

	if updated.NotificationEmail != "test+updated@christippett.dev" {
		t.Error("domain returned unexpected object")
	}

	if err = c.DeleteDomain(ctx, updated); err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_AliasCRUD(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()
	d := "example.com"

	_, err := c.AddDomain(ctx, &Domain{Domain: d})
	if err != nil {
		t.Fatal(err)
	}

	alias, err := c.CreateAlias(ctx, d, &Alias{
		Alias:   "test",
		Forward: "test@christippett.dev",
	})
	if err != nil {
		t.Fatal(err)
	}

	alias.Forward = "test+updated@christippett.dev"
	alias, err = c.UpdateAlias(ctx, d, alias)
	if err != nil {
		t.Fatal(err)
	}

	domain, err := c.GetDomain(ctx, d)
	if err != nil {
		t.Fatal(err)
	}

	defer c.DeleteDomain(ctx, domain)

	// +1 alias to account for the default alias '*'
	aliasCount := len(*domain.Aliases)
	if aliasCount != 2 {
		t.Errorf("domain has unexpected alias count: %d", aliasCount)
	}

	// compare updated alias with last alias
	a := (*domain.Aliases)[len(*domain.Aliases)-1]
	if a.Alias != alias.Alias || a.Forward != alias.Forward {
		t.Error("updated alias does not match domain alias")
	}
}

func TestIntegration_CredentialsCRUD(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()
	d := "example.com"

	domain, err := c.AddDomain(ctx, &Domain{Domain: d})
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.CreateSMTPCredential(ctx, d, &WriteSMTPCredential{
		Username: "test-user",
		Password: "password123",
	})
	if err != nil {
		t.Fatal(err)
	}

	creds, err := c.ListSMTPCredentials(ctx, d)
	if err != nil {
		t.Fatal(err)
	}

	defer c.DeleteDomain(ctx, domain)

	credsCount := len(*creds)
	if credsCount != 1 {
		t.Errorf("domain has unexpected credential count: %d", credsCount)
	}

	// compare updated alias with last alias
	created := (*creds)[credsCount-1]
	if created.Username != "test-user" {
		t.Error("created credential does not match domain")
	}
}
