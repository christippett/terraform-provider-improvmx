package improvmx

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDNS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "improvmx_domain" "test" {
					domain = "%[1]s"
				}

				data "improvmx_dns" "test" {
					domain = improvmx_domain.test.domain
				}
				`, testDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs(
						"data.improvmx_dns.test",
						"records.*",
						map[string]*regexp.Regexp{"type": regexp.MustCompile("MX|TXT|CNAME")},
					),
				),
			},
		},
	})
}
