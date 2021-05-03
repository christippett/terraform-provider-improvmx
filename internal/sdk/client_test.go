package improvmx

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupClient(t *testing.T) Client {
	apiKey := os.Getenv("IMPROVMX_API_KEY")
	if apiKey == "" {
		t.Fatal("'IMPROVMX_API_KEY' must be set for tests")
	}
	return NewClient("https://api.improvmx.com/v3", apiKey, nil)
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
	var count int

	// get current domain count
	domains, err := c.ListDomains(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	initCount := len(*domains)

	// +1 domain
	domain, err := c.AddDomain(ctx, &Domain{Domain: "example.com"})
	if err != nil {
		t.Fatal(err)
	}

	domainsAfterCreate, err := c.ListDomains(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	count = len(*domainsAfterCreate)

	if count != initCount+1 {
		t.Errorf("unexpected domain count: wanted %d, got %d", initCount+1, count)
	}

	// -1 domain
	if err = c.DeleteDomain(ctx, domain); err != nil {
		t.Fatal(err)
	}

	domainsAfterDelete, err := c.ListDomains(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	count = len(*domainsAfterDelete)

	if initCount != count {
		t.Errorf("unexpected domain count: wanted %d, got %d", initCount, count)
	}
}

func TestIntegration_UpdateDomain(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	domain, err := c.AddDomain(ctx, &Domain{Domain: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteDomain(ctx, domain)

	webhook := "http://example.com/webhook"
	domain.Webhook = webhook
	_, err = c.UpdateDomain(ctx, domain)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := c.GetDomain(ctx, "example.com")
	if err != nil {
		t.Fatal(err)
	}

	if updated.Webhook != webhook {
		t.Errorf(
			"domain returned unexpected value for 'webhook': wanted %s, got %s",
			webhook,
			updated.Webhook,
		)
	}
}

func TestIntegration_CheckDomain(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	domain, err := c.AddDomain(ctx, &Domain{Domain: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteDomain(ctx, domain)

	check, err := c.CheckDomain(ctx, domain)
	if err != nil {
		t.Fatal(err)
	}

	// check MX record (value slice)
	expected := RecordValues{
		"mx1.improvmx.com",
		"mx2.improvmx.com",
	}
	if !cmp.Equal(*check.Mx.Expected, expected) {
		t.Errorf("domain check returned unexpected MX record: wanted %s, got %s", expected, check.Mx.Expected)
	}

	// check SPF record (single value)
	expected = RecordValues{"v=spf1 include:spf.improvmx.com -all"}
	if !cmp.Equal(*check.Spf.Expected, expected) {
		t.Errorf("domain check returned unexpected SPF record: wanted %s, got %s", expected, check.Spf.Expected)
	}
}

func TestIntegration_Alias(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()
	d := "example.com"

	domain, err := c.AddDomain(ctx, &Domain{Domain: d})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteDomain(ctx, domain)

	alias, err := c.CreateAlias(ctx, d, &Alias{
		Alias:   "test",
		Forward: "test@piedpiper.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	alias.Forward = "test+updated@piedpiper.com"
	alias, err = c.UpdateAlias(ctx, d, alias)
	if err != nil {
		t.Fatal(err)
	}

	aliases, err := c.ListAliases(ctx, d)
	if err != nil {
		t.Fatal(err)
	}

	// +1 alias to account for the default alias '*'
	aliasCount := len(*aliases)
	if aliasCount != 2 {
		t.Errorf("domain has unexpected alias count: %d", aliasCount)
	}

	// compare updated alias with last alias
	a := (*aliases)[aliasCount-1]
	if a.Alias != alias.Alias || a.Forward != alias.Forward {
		t.Error("updated alias does not match domain alias")
	}

	if err = c.DeleteAlias(ctx, d, alias); err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_SMTPCredentials(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()
	d := "example.com"

	domain, err := c.AddDomain(ctx, &Domain{Domain: d})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteDomain(ctx, domain)

	credential, err := c.CreateSMTPCredential(ctx, d, &WriteSMTPCredential{
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

	credsCount := len(*creds)
	if credsCount != 1 {
		t.Errorf("domain has unexpected credential count: %d", credsCount)
	}

	// compare updated alias with last alias
	created := (*creds)[credsCount-1]
	if created.Username != "test-user" {
		t.Error("created credential does not match domain")
	}

	if err = c.DeleteSMTPCredential(ctx, d, credential); err != nil {
		t.Fatal(err)
	}
}
