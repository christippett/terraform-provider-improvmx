package improvmx

import (
	"context"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomainCheck() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Data source for checking if the MX entries on a domain are valid.",

		ReadContext: dataSourceDomainCheckRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Name of the domain.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"provider_name": {
				Description: "Name of the domain's provider.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"advanced": {
				Description: "True if domain includes advanced configuration.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"valid": {
				Description: "True if domain is currently valid.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"mx":    recordSchema("MX records."),
			"dkim1": recordSchema("DKIM records."),
			"dkim2": recordSchema("DKIM records."),
			"dmarc": recordSchema("DMARC records."),
			"spf":   recordSchema("SPF records."),
		},
	}
}

func recordSchema(description string) *schema.Schema {
	return &schema.Schema{
		Description: description,
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expected": {
					Description: "Expected records for the domain.",
					Type:        schema.TypeSet,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"values": {
					Description: "Current records on the domain.",
					Type:        schema.TypeSet,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"valid": {
					Description: "True if the records on the domain match the expected values.",
					Type:        schema.TypeBool,
					Computed:    true,
				},
			},
		},
	}
}

func dataSourceDomainCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)
	id := d.Get("domain").(string)

	check, err := c.CheckDomain(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	if err = checkResourceData(check, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func checkResourceData(check *improvmx.Check, d *schema.ResourceData) error {
	d.Set("provider_name", check.Provider)
	d.Set("advanced", check.Advanced)
	d.Set("valid", check.Valid)
	d.Set("mx", recordResourceData(check.Mx))
	d.Set("dkim1", recordResourceData(check.Dkim1))
	d.Set("dkim2", recordResourceData(check.Dkim2))
	d.Set("dmarc", recordResourceData(check.Dmarc))
	d.Set("spf", recordResourceData(check.Spf))
	return nil
}

func recordResourceData(record *improvmx.Record) *[]map[string]interface{} {
	m := make([]map[string]interface{}, 1)
	m[0] = map[string]interface{}{
		"expected": recordSet(record.Expected),
		"values":   recordSet(record.Values),
		"valid":    record.Valid,
	}

	return &m
}

func recordSet(rv *improvmx.RecordValues) *schema.Set {
	if rv == nil {
		return nil
	}
	v := *rv
	return schema.NewSet(schema.HashString, v.Interface())
}
