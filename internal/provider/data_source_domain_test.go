package improvmx

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "improvmx_domain" "test" {
					domain = "%[1]s"

					alias {
						alias = "hello"
						forward = "hello@piedpiper.com"
					}
				}

				data "improvmx_domain" "test" {
					domain = "%[1]s"
					depends_on = [improvmx_domain.test]
				}
				`, testDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.improvmx_domain.test", "domain",
						testDomain,
					),
					resource.TestMatchTypeSetElemNestedAttrs(
						"data.improvmx_domain.test",
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
