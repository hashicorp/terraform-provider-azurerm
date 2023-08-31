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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions"
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

func resourceKustoEventGridDataConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoEventGridDataConnectionCreateUpdate,
		Update: resourceKustoEventGridDataConnectionCreateUpdate,
		Read:   resourceKustoEventGridDataConnectionRead,
		Delete: resourceKustoEventGridDataConnectionDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseDataConnectionEventGridV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := dataconnections.ParseDataConnectionID(id)
			return err
		}, importDataConnection("EventGrid")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
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

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"eventhub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubs.ValidateEventhubID,
			},

			"eventhub_consumer_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
			},

			"blob_storage_event_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(dataconnections.BlobStorageEventTypeMicrosoftPointStoragePointBlobCreated),
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForBlobStorageEventType(), false),
			},

			"skip_first_record": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"mapping_rule_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"data_format": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForEventGridDataFormat(), false),
			},

			"database_routing_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      string(dataconnections.DatabaseRoutingSingle),
				ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
			},

			// TODO: rename this to `eventgrid_event_subscription_id` in 4.0
			"eventgrid_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: eventsubscriptions.ValidateScopedEventSubscriptionID,
			},

			// TODO: rename this to `managed_identity_id` in 4.0
			"managed_identity_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.Any(
					clusters.ValidateClusterID,
					commonids.ValidateUserAssignedIdentityID,
				),
			},
		},
	}
}

func resourceKustoEventGridDataConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Event Grid Data Connection creation.")

	id := dataconnections.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_eventgrid_data_connection", id.ID())
		}
	}

	dataConnection := dataconnections.EventGridDataConnection{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &dataconnections.EventGridConnectionProperties{
			StorageAccountResourceId: d.Get("storage_account_id").(string),
			EventHubResourceId:       d.Get("eventhub_id").(string),
			ConsumerGroup:            d.Get("eventhub_consumer_group_name").(string),
			IgnoreFirstRecord:        utils.Bool(d.Get("skip_first_record").(bool)),
		},
	}

	blobStorageEventType := dataconnections.BlobStorageEventType(d.Get("blob_storage_event_type").(string))
	dataConnection.Properties.BlobStorageEventType = &blobStorageEventType

	if tableName, ok := d.GetOk("table_name"); ok {
		dataConnection.Properties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		dataConnection.Properties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		dataFormat := dataconnections.EventGridDataFormat(df.(string))
		dataConnection.Properties.DataFormat = &dataFormat
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		databaseRoutingType := dataconnections.DatabaseRouting(databaseRouting.(string))
		dataConnection.Properties.DatabaseRouting = &databaseRoutingType
	}

	if eventGridRID, ok := d.GetOk("eventgrid_resource_id"); ok {
		dataConnection.Properties.EventGridResourceId = utils.String(eventGridRID.(string))
	}

	if managedIdentityRID, ok := d.GetOk("managed_identity_resource_id"); ok {
		dataConnection.Properties.ManagedIdentityResourceId = utils.String(managedIdentityRID.(string))
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, dataConnection)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceKustoEventGridDataConnectionRead(d, meta)
}

func resourceKustoEventGridDataConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		dataConnection := (*resp.Model).(dataconnections.EventGridDataConnection)
		d.Set("location", location.NormalizeNilable(dataConnection.Location))
		if props := dataConnection.Properties; props != nil {
			d.Set("storage_account_id", props.StorageAccountResourceId)
			d.Set("eventhub_id", props.EventHubResourceId)
			d.Set("eventhub_consumer_group_name", props.ConsumerGroup)
			d.Set("skip_first_record", props.IgnoreFirstRecord)
			d.Set("blob_storage_event_type", string(pointer.From(props.BlobStorageEventType)))
			d.Set("table_name", props.TableName)
			d.Set("mapping_rule_name", props.MappingRuleName)
			d.Set("data_format", string(pointer.From(props.DataFormat)))
			d.Set("database_routing_type", string(pointer.From(props.DatabaseRouting)))
			d.Set("eventgrid_resource_id", props.EventGridResourceId)

			managedIdentityResourceId := ""
			if props.ManagedIdentityResourceId != nil && *props.ManagedIdentityResourceId != "" {
				managedIdentityResourceId = *props.ManagedIdentityResourceId
				clusterId, clusterIdErr := clusters.ParseClusterIDInsensitively(managedIdentityResourceId)
				if clusterIdErr == nil {
					managedIdentityResourceId = clusterId.ID()
				} else {
					userAssignedIdentityId, userAssignedIdentityIdErr := commonids.ParseUserAssignedIdentityIDInsensitively(managedIdentityResourceId)
					if userAssignedIdentityIdErr == nil {
						managedIdentityResourceId = userAssignedIdentityId.ID()
					} else {
						return fmt.Errorf("parsing `managed_identity_resource_id`: %+v; %+v", clusterIdErr, userAssignedIdentityIdErr)
					}
				}
			}
			d.Set("managed_identity_resource_id", managedIdentityResourceId)
		}
	}

	return nil
}

func resourceKustoEventGridDataConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
