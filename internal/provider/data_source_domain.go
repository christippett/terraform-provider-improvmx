package improvmx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Data source for domains.",

		ReadContext: dataSourceDomainRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Name of the domain.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"notification_email": {
				Description: "Email to send the notifications to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"whitelabel": {
				Description: "Parent’s domain that will be displayed for the DNS settings.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active": {
				Description: "True if domain is currently active.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"display": {
				Description: "Display name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dkim_selector": {
				Description: "DKIM selector for domain.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"webhook": {
				Description: "Endpoint to send email events to as POST requests.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"added": {
				Description: "Timestamp when the domain was added.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"dns": {
				Description: "Domain DNS records.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Resource record type. Example: `MX`. Possible values are `MX`, `TXT`, and `CNAME`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Relative name of the object affected by this record. Only applicable for CNAME records. Example: 'dkimprovmx1._domainkey'.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"value": {
							Description: "Data for this record.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"alias": {
				Description: "List of domain alias.",
				Type:        schema.TypeSet,
				Set:         hashSetValue("alias"),
				MinItems:    0,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Description: "Alias to be used in front of your domain, like “contact”, “info”, etc.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"forward": {
							Description: "Destination email to forward the emails to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "Unique ID for alias.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("domain").(string))
	return resourceDomainRead(ctx, d, meta)
}
