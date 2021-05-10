package improvmx

import (
	"context"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNS() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Data source containing the DNS records needed for a domain to work with ImprovMX.",

		ReadContext: dataSourceDNSRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Name of the domain.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"records": {
				Description: "List of domain alias.",
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
		},
	}
}

func dataSourceDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)
	id := d.Get("domain").(string)

	check, err := c.CheckDomain(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	var records []map[string]interface{}
	records = append(records, makeRecords(check.Mx.Expected, "MX", "")...)
	records = append(records, makeRecords(check.Spf.Expected, "TXT", "")...)
	records = append(records, makeRecords(check.Dmarc.Expected, "TXT", "")...)
	records = append(records, makeRecords(check.Dkim1.Expected, "CNAME", "dkimprovmx1._domainkey")...)
	records = append(records, makeRecords(check.Dkim2.Expected, "CNAME", "dkimprovmx2._domainkey")...)

	d.SetId(id)
	d.Set("records", records)

	return nil
}

func makeRecords(r *improvmx.RecordValues, recordType, name string) []map[string]interface{} {
	if r == nil {
		return nil
	}
	rec := make([]map[string]interface{}, len(*r))
	for i, v := range *r {
		rec[i] = map[string]interface{}{
			"type":  recordType,
			"name":  name,
			"value": v,
		}
	}
	return rec
}
