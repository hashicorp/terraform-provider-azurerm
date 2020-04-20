package iothub

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceIotHubDPSSharedAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIotHubDPSSharedAccessPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9!._-]{1,64}`), ""+
					"The shared access policy key name must not be empty, and must not exceed 64 characters in length.  The shared access policy key name can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens."),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"iothub_dps_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"primary_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceIotHubDPSSharedAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	iothubDpsName := d.Get("iothub_dps_name").(string)

	iothubDps, err := client.Get(ctx, iothubDpsName, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(iothubDps.Response) {
			return fmt.Errorf("Error: IotHub DPS %q (Resource Group %q) was not found", iothubDpsName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving IotHub DPS %q (Resource Group %q): %+v", iothubDpsName, resourceGroup, err)
	}

	accessPolicy, err := client.ListKeysForKeyName(ctx, iothubDpsName, keyName, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			return fmt.Errorf("Error: Shared Access Policy %q (IotHub DPS %q / Resource Group %q) was not found", keyName, iothubDpsName, resourceGroup)
		}

		return fmt.Errorf("Error loading Shared Access Policy %q (IotHub DPS %q / Resource Group %q): %+v", keyName, iothubDpsName, resourceGroup, err)
	}

	d.Set("name", keyName)
	d.Set("resource_group_name", resourceGroup)

	resourceID := fmt.Sprintf("%s/keys/%s", *iothubDps.ID, keyName)
	d.SetId(resourceID)

	d.Set("primary_key", accessPolicy.PrimaryKey)
	d.Set("secondary_key", accessPolicy.SecondaryKey)

	primaryConnectionString := ""
	secondaryConnectionString := ""
	if iothubDps.Properties != nil && iothubDps.Properties.ServiceOperationsHostName != nil {
		hostname := iothubDps.Properties.ServiceOperationsHostName
		if primary := accessPolicy.PrimaryKey; primary != nil {
			primaryConnectionString = getSAPConnectionString(*hostname, keyName, *primary)
		}
		if secondary := accessPolicy.SecondaryKey; secondary != nil {
			secondaryConnectionString = getSAPConnectionString(*hostname, keyName, *secondary)
		}
	}
	d.Set("primary_connection_string", primaryConnectionString)
	d.Set("secondary_connection_string", secondaryConnectionString)

	return nil
}
