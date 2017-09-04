package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCreateUpdate,
		Read:   resourceArmAppServiceRead,
		Update: resourceArmAppServiceCreateUpdate,
		Delete: resourceArmAppServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"skip_dns_registration": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_custom_domain_verification": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"force_dns_registration": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ttl_in_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"site_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_service_plan_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"always_on": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"location": locationSchema(),
			"tags":     tagsSchema(),
			"delete_metrics": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceArmAppServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(*ArmClient).appsClient

	log.Printf("[INFO] preparing arguments for Azure ARM App Service creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	skipDNSRegistration := d.Get("skip_dns_registration").(bool)
	skipCustomDomainVerification := d.Get("skip_custom_domain_verification").(bool)
	forceDNSRegistration := d.Get("force_dns_registration").(bool)
	ttlInSeconds := 0
	if v, ok := d.GetOk("ttl_in_seconds"); ok {
		ttlInSeconds = v.(int)
	}
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	siteProps := expandAzureRmAppServiceSiteProps(d)

	siteEnvelope := web.Site{
		Location:       &location,
		Tags:           expandTags(tags),
		SiteProperties: siteProps,
	}

	_, error := appClient.CreateOrUpdate(resGroup, name, siteEnvelope, &skipDNSRegistration, &skipCustomDomainVerification, &forceDNSRegistration, strconv.Itoa(ttlInSeconds), make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := appClient.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceRead(d, meta)
}

func resourceArmAppServiceRead(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(*ArmClient).appsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading App Service details %s", id)

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := appClient.Get(resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if siteProps := resp.SiteProperties; siteProps != nil {
		d.Set("site_config", flattenAzureRmAppServiceSiteProps(siteProps))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServiceDelete(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(*ArmClient).appsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	log.Printf("[DEBUG] Deleting App Service %s: %s", resGroup, name)

	deleteMetrics := d.Get("delete_metrics").(bool)
	deleteEmptyServerFarm := true
	skipDNSRegistration := d.Get("skip_dns_registration").(bool)

	_, err = appClient.Delete(resGroup, name, &deleteMetrics, &deleteEmptyServerFarm, &skipDNSRegistration)

	return err
}

func expandAzureRmAppServiceSiteProps(d *schema.ResourceData) *web.SiteProperties {
	configs := d.Get("site_config").([]interface{})
	siteProps := web.SiteProperties{}
	if len(configs) == 0 {
		return &siteProps
	}
	config := configs[0].(map[string]interface{})

	siteConfig := web.SiteConfig{}
	alwaysOn := config["always_on"].(bool)
	siteConfig.AlwaysOn = utils.Bool(alwaysOn)

	siteProps.SiteConfig = &siteConfig

	serverFarmID := config["app_service_plan_id"].(string)
	siteProps.ServerFarmID = &serverFarmID

	return &siteProps
}

func flattenAzureRmAppServiceSiteProps(siteProps *web.SiteProperties) []interface{} {
	result := make([]interface{}, 0, 1)
	site_config := make(map[string]interface{}, 0)

	if siteProps.ServerFarmID != nil {
		site_config["app_service_plan_id"] = *siteProps.ServerFarmID
	}

	siteConfig := siteProps.SiteConfig
	log.Printf("[DEBUG] SiteConfig is %s", siteConfig)
	if siteConfig != nil {
		if siteConfig.AlwaysOn != nil {
			site_config["always_on"] = *siteConfig.AlwaysOn
		}
	}

	result = append(result, site_config)
	return result
}
