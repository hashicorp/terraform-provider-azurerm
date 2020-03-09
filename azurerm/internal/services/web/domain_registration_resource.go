package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDomainRegistration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDomainRegistrationCreate,
		Read:   resourceArmDomainRegistrationRead,
		Update: resourceArmDomainRegistrationUpdate,
		Delete: resourceArmDomainRegistrationDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DomainRegistrationID(id)
			return err
		}),
		// todo - find sensible values
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Read:   schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DomainRegistrationName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"admin_contact": domainRegistrationContactSchema(),

			"billing_contact": domainRegistrationContactSchema(),

			"registrant_contact": domainRegistrationContactSchema(),

			"technical_contact": domainRegistrationContactSchema(),

			// Optional
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"dns_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AzureDNS",
				ValidateFunc: validation.StringInSlice([]string{
					string(web.AzureDNS),
					string(web.DefaultDomainRegistrarDNS),
				}, false),
			},

			"dns_zone_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"privacy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tags.Schema(),

			// Computed
			"consent": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agreed_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agreed_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agreement_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_renewed": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDomainRegistrationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Failure checking for presence of existing Domain Registration %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_domain_registration", *existing.ID)
		}
	}

	location := azure.NormalizeLocation("global")
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	adminContact := expandDomainRegistrationContact(d.Get("admin_contact").([]interface{}))
	billingContact := expandDomainRegistrationContact(d.Get("billing_contact").([]interface{}))
	registrantContact := expandDomainRegistrationContact(d.Get("registrant_contact").([]interface{}))
	techContact := expandDomainRegistrationContact(d.Get("technical_contact").([]interface{}))

	nameParts := strings.Split(name, ".")
	tld := nameParts[len(nameParts)-1]

	consent, err := getDomainPurchaseConsent(ctx, meta, tld, adminContact.Email)
	if err != nil {
		return err
	}

	autoRenew := d.Get("auto_renew").(bool)
	dnsType := d.Get("dns_type").(string)

	domain := web.Domain{
		Name:     &name,
		Location: &location,
		//Kind:             nil,  // TODO - Find documentation on this
		DomainProperties: &web.DomainProperties{
			ContactAdmin:      adminContact,
			ContactBilling:    billingContact,
			ContactRegistrant: registrantContact,
			ContactTech:       techContact,
			Privacy:           utils.Bool(d.Get("privacy").(bool)),
			AutoRenew:         &autoRenew,
			Consent:           consent,
			//DomainNotRenewableReasons: nil,
			DNSType: web.DNSType(dnsType),
			//AuthCode:                    nil, // TODO - Find documentation on this
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("dns_zone_id"); ok {
		dnsZoneID := v.(string)
		domain.DNSZoneID = &dnsZoneID
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, domain)
	if err != nil {
		return fmt.Errorf("Failure in creating Domain Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Failure waiting for creation of Domain Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Failure reading Domain Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("ID for Domain Registration %q (Resource Group %q) was nil", name, resourceGroup)
	}

	return resourceArmDomainRegistrationRead(d, meta)
}

func resourceArmDomainRegistrationUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO - Stubbed for create/delete testing
	return resourceArmDomainRegistrationRead(d, meta)
}

func resourceArmDomainRegistrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainRegistrationID(d.Id())
	if err != nil {
		return err
	}

	name := id.Name
	resourceGroup := id.ResourceGroup

	domain, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading Domain Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azure.NormalizeLocation(domain.Location))
	d.Set("admin_contact", flattenDomainRegistrationContact(domain.ContactAdmin))
	d.Set("billing_contact", flattenDomainRegistrationContact(domain.ContactBilling))
	d.Set("registrant_contact", flattenDomainRegistrationContact(domain.ContactRegistrant))
	d.Set("technical_contact", flattenDomainRegistrationContact(domain.ContactTech))
	d.Set("consent", flattenDomainRegistrationPurchaseConsent(domain.Consent))
	d.Set("auto_renew", domain.AutoRenew)
	d.Set("created_time", domain.CreatedTime.String())
	d.Set("expiration_time", domain.ExpirationTime.String())
	d.Set("last_renewed", domain.LastRenewedTime.String())
	d.Set("kind", domain.Kind)

	return tags.FlattenAndSet(d, domain.Tags)
}

func resourceArmDomainRegistrationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainRegistrationID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	log.Printf("[DEBUG] Deleting Domain Registration %q (Resource Group %q)", name, resourceGroup)

	// deletes immediately.  `false` here means 24h soft-delete then purged automatically
	forceHardDelete := utils.Bool(true)
	resp, err := client.Delete(ctx, resourceGroup, name, forceHardDelete)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error deleting Domain Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}
