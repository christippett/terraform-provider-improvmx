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
			"alias": {
				Description: "List of domain aliases.",
				Type:        schema.TypeSet,
				Set:         hashAlias,
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

func hashAlias(v interface{}) int {
	alias := v.(map[string]interface{})["alias"].(string)
	return hash(alias)
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(improvmx.Client)

	// add domain
	addDomain := improvmx.Domain{
		Domain:            d.Get("domain").(string),
		NotificationEmail: d.Get("notification_email").(string),
		Whitelabel:        d.Get("whitelabel").(string),
		Webhook:           d.Get("webhook").(string),
	}
	domain, err := c.AddDomain(ctx, &addDomain)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(domain.Domain)

	aliases := d.Get("alias").(*schema.Set)
	count := aliases.Len()

	if count > 0 {
		defaultAliases, err := c.ListAliases(ctx, domain.Domain)
		if err != nil {
			return diag.FromErr(err)
		}

		// delete default aliases when aliases are explicitly defined
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

		// add aliases after domain has been created
		for _, a := range aliases.List() {
			aMap := a.(map[string]interface{})
			alias := improvmx.Alias{
				Alias:   aMap["alias"].(string),
				Forward: aMap["forward"].(string),
			}
			_, err = c.CreateAlias(ctx, domain.Domain, &alias)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  err.Error(),
				})
			}
		}
		if diags.HasError() {
			return diags
		}

		// refresh domain before populating resource
		return resourceDomainRead(ctx, d, meta)
	}

	// ignore domain aliases if not defined on the resource
	domain.Aliases = nil
	domainResourceData(domain, d)
	return nil
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)
	domain, err := c.GetDomain(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = domainResourceData(domain, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)
	updateDomain := improvmx.Domain{
		Domain:            d.Id(),
		NotificationEmail: d.Get("notification_email").(string),
		Whitelabel:        d.Get("whitelabel").(string),
		Webhook:           d.Get("webhook").(string),
	}
	domain, err := c.UpdateDomain(ctx, &updateDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = domainResourceData(domain, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(improvmx.Client)
	err := c.DeleteDomain(ctx, &improvmx.Domain{Domain: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func domainResourceData(domain *improvmx.Domain, d *schema.ResourceData) error {
	d.Set("domain", domain.Domain)
	d.Set("active", domain.Active)
	d.Set("display", domain.Display)
	d.Set("dkim_selector", domain.Whitelabel)
	d.Set("notification_email", domain.NotificationEmail)
	d.Set("webhook", domain.Webhook)
	d.Set("whitelabel", domain.Whitelabel)
	d.Set("added", domain.Added)

	aliasList := make([]interface{}, len(*domain.Aliases))
	for i, a := range *domain.Aliases {
		alias := map[string]interface{}{
			"alias":   a.Alias,
			"forward": a.Forward,
			"id":      a.ID,
		}
		aliasList[i] = alias
	}
	aliases := schema.NewSet(hashAlias, aliasList)
	d.Set("alias", aliases)

	return nil
}
