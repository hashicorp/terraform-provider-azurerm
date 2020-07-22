package web

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/validate"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceDomainCreateUpdate,
		Read:   resourceArmAppServiceDomainRead,
		Update: resourceArmAppServiceDomainCreateUpdate,
		Delete: resourceArmAppServiceDomainDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServiceDomainID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"contact": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"primary_address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"city": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"country": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"postal_code": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"state": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"secondary_address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"email": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"first_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"last_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"fax": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"job_title": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"middle_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"organization": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			// TODO: "consent" is required

			"dns_zone_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DnsZoneID,
			},

			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"dns_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.AzureDNS),
					string(web.DefaultDomainRegistrarDNS),
				}, false),
				Default: string(web.AzureDNS),
			},

			"privacy_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"name_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAppServiceDomainCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing App Service Domain %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_domain", *resp.ID)
		}
	}

	contact := expandAppServiceDomainContact(d.Get("contact").([]interface{}))

	props := web.Domain{
		Location: utils.String("global"),
		DomainProperties: &web.DomainProperties{
			ContactAdmin:      contact,
			ContactBilling:    contact,
			ContactRegistrant: contact,
			ContactTech:       contact,
			Consent:           &web.DomainPurchaseConsent{},
			AutoRenew:         utils.Bool(d.Get("auto_renew").(bool)),
			DNSType:           web.DNSType(d.Get("dns_type").(string)),
			Privacy:           utils.Bool(d.Get("privacy_enabled").(bool)),
			DNSZoneID:         utils.String(d.Get("dns_zone_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, props)
	if err != nil {
		return fmt.Errorf("creating App Service Domain %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of App Service Domain %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// TODO: determine whether need to poll registrationStatus

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving App Service Domain %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for App Service Domain %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmAppServiceDomainRead(d, meta)
}

func resourceArmAppServiceDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Domain %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving App Service Domain %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if props := resp.DomainProperties; props != nil {
		// The Contact information is absent in response, hence we do not flatten "contact" here.
		d.Set("auto_renew", props.AutoRenew)
		d.Set("dns_type", string(props.DNSType))
		d.Set("dns_zone_id", props.DNSZoneID)
		d.Set("privacy_enabled", props.Privacy)
		if err := d.Set("name_servers", utils.FlattenStringSlice(props.NameServers)); err != nil {
			return fmt.Errorf("settting `name_servers`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.DomainClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceDomainID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name, utils.Bool(true)); err != nil {
		return fmt.Errorf("deleting App Service Domain %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandAppServiceDomainContact(input []interface{}) *web.Contact {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &web.Contact{
		AddressMailing: expandAppServiceDomainAddress(raw["address"].([]interface{})),
		Email:          utils.String(raw["email"].(string)),
		NameFirst:      utils.String(raw["first_name"].(string)),
		NameMiddle:     utils.String(raw["middle_name"].(string)),
		NameLast:       utils.String(raw["last_name"].(string)),
		Phone:          utils.String(raw["phone"].(string)),
		Fax:            utils.String(raw["fax"].(string)),
		JobTitle:       utils.String(raw["job_title"].(string)),
		Organization:   utils.String(raw["organization"].(string)),
	}

	return output
}

func expandAppServiceDomainAddress(input []interface{}) *web.Address {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &web.Address{
		Address1:   utils.String(raw["primary_address"].(string)),
		Address2:   utils.String(raw["secondary_address"].(string)),
		City:       utils.String(raw["city"].(string)),
		Country:    utils.String(raw["country"].(string)),
		PostalCode: utils.String(raw["postal_code"].(string)),
		State:      utils.String(raw["state"].(string)),
	}
	return output
}
