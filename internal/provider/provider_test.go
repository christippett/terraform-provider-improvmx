package improvmx

import (
	"os"
	"testing"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var improvmxClient improvmx.Client
var providerFactories map[string]func() (*schema.Provider, error)

func init() {
	// providerFactories are used to instantiate a provider during acceptance testing.
	// The factory function will be invoked for every Terraform CLI command executed
	// to create a provider server to which the CLI can reattach.
	improvmxClient = improvmx.NewClient("https://api.improvmx.com/v3", os.Getenv("IMPROVMX_API_KEY"), nil)
	providerFactories = map[string]func() (*schema.Provider, error){
		"improvmx": func() (*schema.Provider, error) {
			return New("dev")(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("IMPROVMX_API_KEY"); err == "" {
		t.Fatal("IMPROVMX_API_KEY must be set for acceptance tests")
	}
}
