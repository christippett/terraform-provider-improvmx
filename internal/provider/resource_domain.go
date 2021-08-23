package improvmx

import (
	"context"

	improvmx "github.com/christippett/terraform-provider-improvmx/internal/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				Optional:    true,
			},
			"whitelabel": {
				Description: "Parent’s domain that will be displayed for the DNS settings.",
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:    true,
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
				Description: "List of domain aliases.",
				Type:        schema.TypeSet,
				Set:         hashSetValue("alias"),
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Description: "Alias to be used in front of your domain, like “contact”, “info”, etc.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"forward": {
							Description: "Destination email to forward the emails to.",
							Type:        schema.TypeString,
							Required:    true,
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

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(improvmx.Client)

	// add domain
	inputDomain := domainFromResourceData(d)
	domain, err := c.AddDomain(ctx, inputDomain)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(domain.Domain)

	// domain email aliases
	domain.Aliases = aliasesFromSet(d.Get("alias").(*schema.Set))
	if domain.Aliases != nil {
		// get aliases created by default when the domain is first created
		defaultAliases, err := c.ListAliases(ctx, domain.Domain)
		if err != nil {
			return diag.FromErr(err)
		}

		// delete default alias(es) if the resource has defined its own resources
		for _, a := range *defaultAliases {
			if err = c.DeleteAlias(ctx, domain.Domain, &a); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  err.Error(),
				})
			}
		}
		if diags.HasError() {
			return diags
		}

		// add aliases after domain has been created and any default aliases have
		// been deleted
		for _, alias := range *domain.Aliases {
			_, err = c.CreateAlias(ctx, domain.Domain, &alias)
			if err != nil {
				diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: err.Error()})
			}
		}
		if diags.HasError() {
			return diags
		}
	}

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)

	inputDomain := domainFromResourceData(d)
	domain, err := c.UpdateDomain(ctx, inputDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	domain.Aliases = aliasesFromSet(d.Get("alias").(*schema.Set))
	if domain.Aliases != nil || d.HasChange("alias") {
		old, new := getSetChange(d, "alias")
		// create if alias in new, but not in old
		for _, a := range *aliasesFromSet(new.Difference(old)) {
			_, err = c.CreateAlias(ctx, domain.Domain, &a)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		// delete if alias in old, but not in new
		for _, a := range *aliasesFromSet(old.Difference(new)) {
			if err = c.DeleteAlias(ctx, domain.Domain, &a); err != nil {
				return diag.FromErr(err)
			}
		}
		// update if alias in both old and new
		for _, a := range *aliasesFromSet(new.Intersection(old)) {
			_, err = c.UpdateAlias(ctx, domain.Domain, &a)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)

	domain, err := c.GetDomain(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// check domain
	check, err := c.CheckDomain(ctx, domain.Domain)
	if err != nil {
		return diag.FromErr(err)
	}

	var dns []map[string]interface{}
	dns = append(dns, makeDNSRecord(check.Mx.Expected, "MX", "")...)
	dns = append(dns, makeDNSRecord(check.Spf.Expected, "TXT", "")...)
	dns = append(dns, makeDNSRecord(check.Dmarc.Expected, "TXT", "")...)
	dns = append(dns, makeDNSRecord(check.Dkim1.Expected, "CNAME", "dkimprovmx1._domainkey")...)
	dns = append(dns, makeDNSRecord(check.Dkim2.Expected, "CNAME", "dkimprovmx2._domainkey")...)
	d.Set("dns", dns)

	inputAliases := d.Get("alias").(*schema.Set)
	if inputAliases.Len() == 0 {
		domain.Aliases = nil
	}

	return resourceDataFromDomain(domain, d)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)
	err := c.DeleteDomain(ctx, &improvmx.Domain{Domain: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDataFromDomain(domain *improvmx.Domain, d *schema.ResourceData) diag.Diagnostics {
	d.Set("domain", domain.Domain)
	d.Set("active", domain.Active)
	d.Set("display", domain.Display)
	d.Set("dkim_selector", domain.Whitelabel)
	d.Set("notification_email", domain.NotificationEmail)
	d.Set("webhook", domain.Webhook)
	d.Set("whitelabel", domain.Whitelabel)
	d.Set("added", domain.Added)

	// return early if there's no aliases to process
	if domain.Aliases == nil {
		return nil
	}

	aliasList := make([]interface{}, len(*domain.Aliases))
	for i, a := range *domain.Aliases {
		aliasList[i] = map[string]interface{}{
			"alias":   a.Alias,
			"forward": a.Forward,
			"id":      a.ID,
		}
	}
	aliases := schema.NewSet(hashSetValue("alias"), aliasList)
	d.Set("alias", aliases)

	dnsList := []map[string]interface{}{
		0: {
			"mx": []string{"mx1.improvmx.com"},
		},
	}
	d.Set("dns_settings", dnsList)

	return nil
}

func aliasesFromSet(s *schema.Set) *[]improvmx.Alias {
	if s.Len() == 0 {
		return nil
	}
	aliases := make([]improvmx.Alias, s.Len())
	for i, a := range s.List() {
		item := a.(map[string]interface{})
		aliases[i] = improvmx.Alias{
			Alias:   item["alias"].(string),
			Forward: item["forward"].(string),
		}
	}
	return &aliases
}

func domainFromResourceData(d *schema.ResourceData) *improvmx.Domain {
	return &improvmx.Domain{
		Domain:            d.Get("domain").(string),
		NotificationEmail: d.Get("notification_email").(string),
		Whitelabel:        d.Get("whitelabel").(string),
		Webhook:           d.Get("webhook").(string),
	}
}

func getSetChange(d *schema.ResourceData, key string) (*schema.Set, *schema.Set) {
	old, new := d.GetChange(key)
	return old.(*schema.Set), new.(*schema.Set)
}

func makeDNSRecord(r *improvmx.RecordValues, recordType, name string) []map[string]interface{} {
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
