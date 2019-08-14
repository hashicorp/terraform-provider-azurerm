package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAppService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAppServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocationForDataSource(),

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"site_config": azure.SchemaAppServiceDataSourceSiteConfig(),

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"client_cert_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"connection_string": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),

			"site_credential": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"default_site_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"possible_outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_control": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceArmAppServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: App Service %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", name, err)
	}

	configResp, err := client.GetConfiguration(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Configuration %q: %+v", name, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service AppSettings %q: %+v", name, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service ConnectionStrings %q: %+v", name, err)
	}

	scmResp, err := client.GetSourceControl(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Source Control %q: %+v", name, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Site Credential %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("enabled", props.Enabled)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("client_cert_enabled", props.ClientCertEnabled)
		d.Set("default_site_hostname", props.DefaultHostName)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
	}

	if err := d.Set("app_settings", flattenAppServiceAppSettings(appSettingsResp.Properties)); err != nil {
		return err
	}
	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	siteConfig := azure.FlattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	scm := flattenAppServiceSourceControl(scmResp.SiteSourceControlProperties)
	if err := d.Set("source_control", scm); err != nil {
		return err
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
