package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceAppServiceSlotConfigNames() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSlotConfigNamesCreateUpdate,
		Read:   resourceAppServiceSlotConfigNamesRead,
		Update: resourceAppServiceSlotConfigNamesCreateUpdate,
		Delete: resourceAppServiceSlotConfigNamesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"app_service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AppServiceName,
			},
			"slot_config_names": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"app_setting_names": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"connection_string_names": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAppServiceSlotConfigNamesCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM App Service slot config names creation.")

	appServiceName := d.Get("app_service_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.ListSlotConfigurationNames(ctx, resourceGroup, appServiceName)
		if err != nil {
			return fmt.Errorf("error querying slot config names settings: %+v", err)
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_slot_config_names", *existing.ID)
		}
	}

	slotConfigNamesResource, err := expandAppServiceSlotConfigNamesResource(d.Get("slot_config_names").(*pluginsdk.Set))
	if err != nil {
		return fmt.Errorf("incorrect slot config names settings: %+v", err)
	}
	newSlotConfigNamesResource, err := client.UpdateSlotConfigurationNames(ctx, resourceGroup, appServiceName, slotConfigNamesResource)
	if err != nil {
		return fmt.Errorf("error creating slot config names for app service %q in resource group %q): %+v", appServiceName, resourceGroup, err)
	}
	d.SetId(*newSlotConfigNamesResource.Name)

	return resourceAppServiceSlotConfigNamesRead(d, meta)
}

func resourceAppServiceSlotConfigNamesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.ListSlotConfigurationNames(ctx, resourceGroup, id)
	if err != nil {
		return fmt.Errorf("error querying slot config names settings: %+v", err)
	}

	d.Set("app_service_name", id)
	d.Set("resource_group_name", resourceGroup)

	if *existing.SlotConfigNames.AppSettingNames != nil && *existing.SlotConfigNames.ConnectionStringNames != nil {
		d.Set("slot_config_names", flattenAppServiceSlotConfigNamesResource(existing))
	}

	return nil
}

func resourceAppServiceSlotConfigNamesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()
	resourceGroup := d.Get("resource_group_name").(string)

	emptyAppSettingNames := make([]string, 0)
	emptyConnectionStringNames := make([]string, 0)

	emptySlotConfigNameResource := web.SlotConfigNamesResource{
		SlotConfigNames: &web.SlotConfigNames{
			AppSettingNames:       &emptyAppSettingNames,
			ConnectionStringNames: &emptyConnectionStringNames,
		},
	}

	_, err := client.UpdateSlotConfigurationNames(ctx, resourceGroup, id, emptySlotConfigNameResource)

	if err != nil {
		return fmt.Errorf("error removing slot config names settings: %+v", err)
	}

	d.SetId("")

	return nil
}

func expandAppServiceSlotConfigNamesResource(input *pluginsdk.Set) (web.SlotConfigNamesResource, error) {
	newSlotConfigNames := input.List()
	appSettingNames := make([]string, 0)
	connectionStringNames := make([]string, 0)
	for _, newSlotConfigName := range newSlotConfigNames {
		slotConfigName := newSlotConfigName.(map[string]interface{})

		for _, v := range slotConfigName["app_setting_names"].([]interface{}) {
			appSettingNames = append(appSettingNames, v.(string))
		}

		for _, v := range slotConfigName["connection_string_names"].([]interface{}) {
			connectionStringNames = append(connectionStringNames, v.(string))
		}
	}

	newAppSettingNames := removeDuplicateValues(appSettingNames)
	newConnectionStringNames := removeDuplicateValues(connectionStringNames)

	slotConfigNameResource := web.SlotConfigNamesResource{
		SlotConfigNames: &web.SlotConfigNames{
			AppSettingNames:       &newAppSettingNames,
			ConnectionStringNames: &newConnectionStringNames,
		},
	}
	return slotConfigNameResource, nil
}

func flattenAppServiceSlotConfigNamesResource(input web.SlotConfigNamesResource) []interface{} {
	results := make([]interface{}, 0)

	result := make(map[string]interface{})
	result["app_setting_names"] = input.SlotConfigNames.AppSettingNames
	result["connection_string_names"] = input.SlotConfigNames.ConnectionStringNames
	results = append(results, result)
	return results
}

func removeDuplicateValues(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
