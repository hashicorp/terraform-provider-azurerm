package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceIotHubDPSSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceIotHubDPSSharedAccessPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IotHubSharedAccessPolicyName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"iothub_dps_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceIotHubDPSSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.DPSResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDpsSharedAccessPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_dps_name").(string), d.Get("name").(string))

	iothubDps, err := client.Get(ctx, id.ProvisioningServiceName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(iothubDps.Response) {
			return fmt.Errorf("Error: IotHub DPS %q (Resource Group %q) was not found", id.ProvisioningServiceName, id.ResourceGroup)
		}

		return fmt.Errorf("retrieving IotHub DPS %q (Resource Group %q): %+v", id.ProvisioningServiceName, id.ResourceGroup, err)
	}

	accessPolicy, err := client.ListKeysForKeyName(ctx, id.ProvisioningServiceName, id.KeyName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.Set("name", id.KeyName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.SetId(id.ID())

	d.Set("primary_key", accessPolicy.PrimaryKey)
	d.Set("secondary_key", accessPolicy.SecondaryKey)

	primaryConnectionString := ""
	secondaryConnectionString := ""
	if iothubDps.Properties != nil && iothubDps.Properties.ServiceOperationsHostName != nil {
		hostname := iothubDps.Properties.ServiceOperationsHostName
		if primary := accessPolicy.PrimaryKey; primary != nil {
			primaryConnectionString = getSAPConnectionString(*hostname, id.KeyName, *primary)
		}
		if secondary := accessPolicy.SecondaryKey; secondary != nil {
			secondaryConnectionString = getSAPConnectionString(*hostname, id.KeyName, *secondary)
		}
	}
	d.Set("primary_connection_string", primaryConnectionString)
	d.Set("secondary_connection_string", secondaryConnectionString)

	return nil
}
