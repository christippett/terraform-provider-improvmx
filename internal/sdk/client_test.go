package improvmx

import (
	"context"
	"os"
	"testing"
	"time"
)

func setupClient(t *testing.T) Client {
	apiKey := os.Getenv("IMPROVMX_API_KEY")
	if apiKey == "" {
		t.Fatal("'IMPROVMX_API_KEY' must be set for tests")
	}
	return NewClient("https://api.improvmx.com/v3", apiKey, nil, os.Stdout)
}

func TestIntegration_ListDomains(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	domain, err := c.AddDomain(ctx, &Domain{
		Domain: "example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	domains, err := c.ListDomains(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer c.DeleteDomain(ctx, domain)

	if len(*domains) < 1 {
		t.Errorf("domain list returned unexpected count")
	}
}

func TestIntegration_GetDomain(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()
	now := time.Now().Unix()
	d := "example.com"

	domain, err := c.AddDomain(ctx, &Domain{
		Domain: d,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.GetDomain(ctx, &d)
	if err != nil {
		t.Fatal(err)
	}

	defer c.DeleteDomain(ctx, domain)

	if got.Domain != d || !(got.Added > now) {
		t.Errorf("domain returned unexpected object: %s", domain.Domain)
	}
}

func TestIntegration_AddDomain(t *testing.T) {
	c := setupClient(t)
	ctx := context.Background()

	domain, err := c.AddDomain(ctx, &Domain{
		Domain:            "example.com",
		NotificationEmail: "test@christippett.dev",
	})
	if err != nil {
		t.Fatal(err)
	}

	defer c.DeleteDomain(ctx, domain)

	if domain.Domain != "example.com" {
		t.Errorf("domain returned unexpected object: %s", domain.Domain)
	}
}

func TestIntegration_UpdateDomain(t *testing.T) {
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
	got, err := c.UpdateDomain(ctx, domain)
	if err != nil {
		t.Fatal(err)
	}

	if got.Domain != "example.com" && got.NotificationEmail == "test+updated@christippett.dev" {
		t.Errorf("domain returned unexpected object: %s vs %s", got.Domain, domain.Domain)
	}

	if err = c.DeleteDomain(ctx, got); err != nil {
		t.Fatal(err)
	}
}
