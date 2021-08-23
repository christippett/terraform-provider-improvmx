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
		Description: "Returns the result of ImprovMX's domain check, including validation of the domain's DNS configuration.",

		ReadContext: dataSourceDomainCheckRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Domain name.",
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
			"mx":    checkRecordSchemaFactory("`mx` record validation."),
			"dkim1": checkRecordSchemaFactory("`dkim1` record validation."),
			"dkim2": checkRecordSchemaFactory("`dkim2` record validation."),
			"dmarc": checkRecordSchemaFactory("`dmarc` record validation."),
			"spf":   checkRecordSchemaFactory("`spf` record validation."),
		},
	}
}

func checkRecordSchemaFactory(description string) *schema.Schema {
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
	d.Set("mx", checkRecordResourceData(check.Mx))
	d.Set("dkim1", checkRecordResourceData(check.Dkim1))
	d.Set("dkim2", checkRecordResourceData(check.Dkim2))
	d.Set("dmarc", checkRecordResourceData(check.Dmarc))
	d.Set("spf", checkRecordResourceData(check.Spf))
	return nil
}

func checkRecordResourceData(record *improvmx.Record) *[]map[string]interface{} {
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
