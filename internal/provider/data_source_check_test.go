package improvmx

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "improvmx_domain" "test" {
					domain = "%[1]s"
				}

				data "improvmx_check" "test" {
					domain = "%[1]s"
					depends_on = [improvmx_domain.test]
				}
				`, testDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.improvmx_check.test",
						"mx.0.expected.0",
						regexp.MustCompile(`mx\d\.improvmx\.com`),
					),
				),
			},
		},
	})
}
