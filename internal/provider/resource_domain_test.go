package improvmx

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDomain(t *testing.T) {
	// t.Skip("resource not yet implemented, remove this once you add your own code")
	domain := "example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(s *terraform.State) error {
			// loop through resources in state and verify each widget has been destroyed
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "improvmx_domain" {
					continue
				}

				domain, err := improvmxClient.GetDomain(context.Background(), rs.Primary.ID)
				if domain != nil && err == nil {
					t.Fatalf("domain '%s' still available after destroy", domain.Domain)
				}
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "improvmx_domain" "test" {
						domain = "%s"

						alias {
							alias = "hello"
							forward = "hello@piedpiper.com"
						}

						alias {
							alias = "contact"
							forward = "contact@piedpiper.com"
						}
					}
				`, domain),
				PreventPostDestroyRefresh: true,
				Destroy:                   false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"improvmx_domain.test",
						"display",
						domain,
					),
					testAccCheckDomainAliasCount("improvmx_domain.test", 2),
					resource.TestMatchTypeSetElemNestedAttrs(
						"improvmx_domain.test",
						"alias.*",
						map[string]*regexp.Regexp{
							"alias":   regexp.MustCompile(`^hello$`),
							"forward": regexp.MustCompile(`^hello\@piedpiper\.com$`),
							"id":      regexp.MustCompile(`\d+`),
						},
					),
				),
			},
		},
	})
}

func testAccCheckDomainAliasCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("domain not set on resource")
		}

		aliases, err := improvmxClient.ListAliases(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		count := len(*aliases)
		if count != expected {
			return fmt.Errorf(
				"unexpected domain alias count: wanted %d, got %d",
				expected, count)
		}
		return nil
	}
}
