package improvmx

import (
	"log"
	"os"
	"testing"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testDomain = "example.com"

var improvmxClient improvmx.Client
var providerFactories map[string]func() (*schema.Provider, error)

func init() {
	improvmxClient = improvmx.NewClient(
		"https://api.improvmx.com/v3",
		os.Getenv("IMPROVMX_API_KEY"),
		log.Writer(),
	)
	providerFactories = map[string]func() (*schema.Provider, error){
		"improvmx": func() (*schema.Provider, error) {
			return New("dev")(), nil
		},
	}
	// resource.AddTestSweepers("improvmx_domain", &resource.Sweeper{
	// 	Name: "improvmx_domain",
	// 	F: func(r string) error {
	// 		ctx := context.Background()
	// 		res, err := improvmxClient.ListDomains(ctx, &improvmx.QueryDomain{
	// 			Query: testDomain,
	// 		})
	// 		if err != nil {
	// 			return fmt.Errorf("error listing domains during test sweep: %s", err)
	// 		}
	// 		if len(*res) == 1 && (*res)[0].Domain == testDomain {
	// 			if err := improvmxClient.DeleteDomain(ctx, &improvmx.Domain{
	// 				Domain: testDomain,
	// 			}); err != nil {
	// 				fmt.Printf("error deleting test domain: %s", err)
	// 			}
	// 		}
	// 		return nil
	// 	},
	// })
}

// func TestMain(m *testing.M) {
// 	resource.TestMain(m)
// }

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
