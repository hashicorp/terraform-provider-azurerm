// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoEventHubDataConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoEventHubDataConnectionCreateUpdate,
		Read:   resourceKustoEventHubDataConnectionRead,
		Update: resourceKustoEventHubDataConnectionCreateUpdate,
		Delete: resourceKustoEventHubDataConnectionDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseDataConnectionEventHubV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := dataconnections.ParseDataConnectionID(id)
			return err
		}, importDataConnection("EventHub")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
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

			"compression": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      dataconnections.CompressionNone,
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForCompression(), false),
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseName,
			},

			"eventhub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubs.ValidateEventhubID,
			},

			"event_system_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"consumer_group": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					eventhubValidate.ValidateEventHubConsumerName(),
					validation.StringInSlice([]string{"$Default"}, false)),
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"identity_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.Any(
					clusters.ValidateClusterID,
					commonids.ValidateUserAssignedIdentityID,
				),
			},

			"mapping_rule_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"data_format": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForEventHubDataFormat(), false),
			},

			"database_routing_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      string(dataconnections.DatabaseRoutingSingle),
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
			},
		},
	}
}

func resourceKustoEventHubDataConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Event Hub Data Connection creation.")

	id := dataconnections.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_eventhub_data_connection", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	eventHubDataConnectionProperties := expandKustoEventHubDataConnectionProperties(d)

	dataConnection1 := dataconnections.EventHubDataConnection{
		Name:       &id.DataConnectionName,
		Location:   &location,
		Properties: eventHubDataConnectionProperties,
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		dbRouting := dataconnections.DatabaseRouting(databaseRouting.(string))
		dataConnection1.Properties.DatabaseRouting = &dbRouting
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, dataConnection1)
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceKustoEventHubDataConnectionRead(d, meta)
}

func resourceKustoEventHubDataConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.DataConnectionName, id.ResourceGroupName, id.ClusterName, id.DatabaseName, err)
	}

	d.Set("name", id.DataConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

	if resp.Model != nil {
		if dataConnection, ok := (*resp.Model).(dataconnections.EventHubDataConnection); ok {
			if location := dataConnection.Location; location != nil {
				d.Set("location", azure.NormalizeLocation(*location))
			}

			if props := dataConnection.Properties; props != nil {
				d.Set("eventhub_id", props.EventHubResourceId)
				d.Set("consumer_group", props.ConsumerGroup)
				d.Set("table_name", props.TableName)
				d.Set("mapping_rule_name", props.MappingRuleName)
				d.Set("data_format", string(pointer.From(props.DataFormat)))
				d.Set("database_routing_type", string(pointer.From(props.DatabaseRouting)))
				d.Set("compression", string(pointer.From(props.Compression)))
				d.Set("event_system_properties", props.EventSystemProperties)

				identityId := ""
				if props.ManagedIdentityResourceId != nil {
					identityId = *props.ManagedIdentityResourceId
					clusterId, clusterIdErr := clusters.ParseClusterIDInsensitively(identityId)
					if clusterIdErr == nil {
						identityId = clusterId.ID()
					} else {
						userAssignedIdentityId, userAssignedIdentityIdErr := commonids.ParseUserAssignedIdentityIDInsensitively(identityId)
						if userAssignedIdentityIdErr == nil {
							identityId = userAssignedIdentityId.ID()
						} else {
							return fmt.Errorf("parsing `identity_id`: %+v; %+v", clusterIdErr, userAssignedIdentityIdErr)
						}
					}
				}
				d.Set("identity_id", identityId)
			}
		}
	}

	return nil
}

func resourceKustoEventHubDataConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnections.ParseDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.DataConnectionName, id.ResourceGroupName, id.ClusterName, id.DatabaseName, err)
	}

	return nil
}

func expandKustoEventHubDataConnectionProperties(d *pluginsdk.ResourceData) *dataconnections.EventHubConnectionProperties {
	eventHubConnectionProperties := &dataconnections.EventHubConnectionProperties{}

	if eventhubResourceID, ok := d.GetOk("eventhub_id"); ok {
		eventHubConnectionProperties.EventHubResourceId = eventhubResourceID.(string)
	}

	if consumerGroup, ok := d.GetOk("consumer_group"); ok {
		eventHubConnectionProperties.ConsumerGroup = consumerGroup.(string)
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		eventHubConnectionProperties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		eventHubConnectionProperties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		dataFormat := dataconnections.EventHubDataFormat(df.(string))
		eventHubConnectionProperties.DataFormat = &dataFormat
	}

	if compression, ok := d.GetOk("compression"); ok {
		comp := dataconnections.Compression(compression.(string))
		eventHubConnectionProperties.Compression = &comp
	}

	if eventSystemProperties, ok := d.GetOk("event_system_properties"); ok {
		props := make([]string, 0)
		for _, prop := range eventSystemProperties.([]interface{}) {
			props = append(props, prop.(string))
		}
		eventHubConnectionProperties.EventSystemProperties = &props
	}

	if identityId, ok := d.GetOk("identity_id"); ok {
		eventHubConnectionProperties.ManagedIdentityResourceId = utils.String(identityId.(string))
	}

	return eventHubConnectionProperties
}
