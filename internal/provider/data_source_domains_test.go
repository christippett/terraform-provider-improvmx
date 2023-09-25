package improvmx

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDomains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "improvmx_domains" "test" { }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs(
						"data.improvmx_domains.test",
						"domains.*",
						map[string]*regexp.Regexp{
							"domain": regexp.MustCompile(`^[a-zA-Z0-9]+\.[a-zA-z]`),
						},
					),
				),
			},
		},
	})
}
