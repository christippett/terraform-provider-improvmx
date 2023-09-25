package improvmx

import (
	"context"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomains() *schema.Resource {
	return &schema.Resource{
		Description: "ImprovMX domain data source.",
		ReadContext: dataSourceDomainsRead,
		Schema: map[string]*schema.Schema{
			"query": {
				Description: "Search the domains starting by this value.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domains": {
				Description: "A list of domain objects linked to an account.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: domainSchema,
				},
			},
		},
	}
}

func dataSourceDomainsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)

	domains, err := c.ListDomains(ctx, &improvmx.QueryDomain{Query: "", IsActive: true})
	if err != nil {
		return diag.FromErr(err)
	}

	domainNames := make([]string, len(*domains))
	domainList := make([]interface{}, len(*domains))
	for i, domain := range *domains {
		domainNames[i] = domain.Domain
		check, err := c.CheckDomain(ctx, domain.Domain)
		if err != nil {
			return diag.FromErr(err)
		}
		domainList[i] = map[string]interface{}{
			"active":             domain.Active,
			"added":              domain.Added,
			"display":            domain.Display,
			"dkim_selector":      domain.DkimSelector,
			"dns":                domainConfigFromCheck(check),
			"domain":             domain.Domain,
			"notification_email": domain.NotificationEmail,
			"webhook":            domain.Webhook,
			"whitelabel":         domain.Whitelabel,
		}
	}
	d.SetId(stringListChecksum(domainNames))
	d.Set("domains", domainList)
	return nil
}
