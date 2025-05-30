// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	iotHubParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoIotHubDataConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoIotHubDataConnectionCreate,
		Read:   resourceKustoIotHubDataConnectionRead,
		Delete: resourceKustoIotHubDataConnectionDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseDataConnectionIotHub0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := dataconnections.ParseDataConnectionID(id)
			return err
		}, importDataConnection("IotHub")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataConnectionName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseName,
			},

			"iothub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IotHubID,
			},

			"consumer_group": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IoTHubConsumerGroupName,
			},

			"shared_access_policy_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IotHubSharedAccessPolicyName,
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"mapping_rule_name": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"data_format": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForIotHubDataFormat(), false),
			},

			"database_routing_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      string(dataconnections.DatabaseRoutingSingle),
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
			},

			"event_system_properties": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceKustoIotHubDataConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Iot Hub Data Connection creation.")

	id := dataconnections.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_kusto_iothub_data_connection", id.ID())
	}

	iotHubDataConnectionProperties := expandKustoIotHubDataConnectionProperties(d)

	dataConnection := dataconnections.IotHubDataConnection{
		Location:   utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: iotHubDataConnectionProperties,
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		dbRoutingType := dataconnections.DatabaseRouting(databaseRouting.(string))
		dataConnection.Properties.DatabaseRouting = &dbRoutingType
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, dataConnection)
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoIotHubDataConnectionRead(d, meta)
}

func resourceKustoIotHubDataConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnections.ParseDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DataConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

	if resp.Model != nil {
		if dataConnection, ok := resp.Model.(dataconnections.IotHubDataConnection); ok {
			d.Set("location", location.NormalizeNilable(dataConnection.Location))
			if props := dataConnection.Properties; props != nil {
				iotHubId := ""
				if parsedIoTHubId, err := iotHubParse.IotHubIDInsensitively(props.IotHubResourceId); err == nil {
					iotHubId = parsedIoTHubId.ID()
				}
				d.Set("iothub_id", iotHubId)
				d.Set("consumer_group", props.ConsumerGroup)
				d.Set("table_name", props.TableName)
				d.Set("mapping_rule_name", props.MappingRuleName)
				d.Set("data_format", string(pointer.From(props.DataFormat)))
				d.Set("database_routing_type", string(pointer.From(props.DatabaseRouting)))
				d.Set("shared_access_policy_name", props.SharedAccessPolicyName)
				d.Set("event_system_properties", utils.FlattenStringSlice(props.EventSystemProperties))
			}
		}
	}

	return nil
}

func resourceKustoIotHubDataConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnections.ParseDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandKustoIotHubDataConnectionProperties(d *pluginsdk.ResourceData) *dataconnections.IotHubConnectionProperties {
	iotHubDataConnectionProperties := &dataconnections.IotHubConnectionProperties{
		IotHubResourceId:       d.Get("iothub_id").(string),
		ConsumerGroup:          d.Get("consumer_group").(string),
		SharedAccessPolicyName: d.Get("shared_access_policy_name").(string),
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		iotHubDataConnectionProperties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		iotHubDataConnectionProperties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		dataFormat := dataconnections.IotHubDataFormat(df.(string))
		iotHubDataConnectionProperties.DataFormat = &dataFormat
	}

	if eventSystemProperties, ok := d.GetOk("event_system_properties"); ok {
		iotHubDataConnectionProperties.EventSystemProperties = utils.ExpandStringSlice(eventSystemProperties.(*pluginsdk.Set).List())
	}

	return iotHubDataConnectionProperties
}
