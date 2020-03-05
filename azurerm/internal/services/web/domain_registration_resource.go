package web

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
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
			Create:  schema.DefaultTimeout(30 * time.Minute),
			Read:    schema.DefaultTimeout(30 * time.Minute),
			Update:  schema.DefaultTimeout(30 * time.Minute),
			Delete:  schema.DefaultTimeout(30 * time.Minute),
			Default: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DomainRegistrationName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"admin_contact": domainRegistrationContactSchema(),

			"billing_contact": domainRegistrationContactSchema(),

			"registrant_contact": domainRegistrationContactSchema(),

			"technical_contact": domainRegistrationContactSchema(),

			"consent": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"agreement_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
						},
						"agreed_by": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"agreed_at": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
					},
				},
			},

			// Optional
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"dns_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.AzureDNS),
					string(web.DefaultDomainRegistrarDNS),
				}, false),
			},

			"dns_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"privacy": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),

			// Computed
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


func resourceArmDomainRegistrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainRegistrationID(d.Id())

	name := id.Name
	resourceGroup := id.ResourceGroup

	domain, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		fmt.Errorf("Error reading Domain Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//Required
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azure.NormalizeLocation(domain.Location))

	d.Set("admin_contact", flattenDomainRegistrationContact(domain.ContactAdmin))
	d.Set("billing_contact", flattenDomainRegistrationContact(domain.ContactBilling))
	d.Set("registrant_contact", flattenDomainRegistrationContact(domain.ContactRegistrant))
	d.Set("technical_contact", flattenDomainRegistrationContact(domain.ContactTech))

	// Computed
	d.Set("created_time", domain.CreatedTime)
	d.Set("expiration_time", domain.ExpirationTime)
	d.Set("last_renewed", domain.LastRenewedTime)
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
