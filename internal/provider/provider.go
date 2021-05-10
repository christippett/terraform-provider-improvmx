package improvmx

import (
	"context"
	"log"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"base_url": {
					Type:     schema.TypeString,
					Optional: true,
					DefaultFunc: schema.EnvDefaultFunc(
						"IMPROVMX_BASE_URL", "https://api.improvmx.com/v3"),
				},
				"api_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("IMPROVMX_API_KEY", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"improvmx_domain": dataSourceDomain(),
				"improvmx_check":  dataSourceDomainCheck(),
				"improvmx_dns":    dataSourceDNS(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"improvmx_domain": resourceDomain(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		url := d.Get("base_url").(string)
		apiKey := d.Get("api_key").(string)
		userAgent := p.UserAgent("terraform-provider-improvmx", version)
		client := improvmx.NewClient(url, apiKey, log.Writer())
		if err := client.SetUserAgent(userAgent); err != nil {
			diag.FromErr(err)
		}
		return client, nil
	}
}
