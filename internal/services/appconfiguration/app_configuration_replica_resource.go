package appconfiguration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ReplicaResource{}

type ReplicaResource struct{}

func (r ReplicaResource) ModelObject() interface{} {
	return &ReplicaResourceSchema{}
}

type ReplicaResourceSchema struct {
	ConfigurationStoreId string `tfschema:"configuration_store_id"`
	Endpoint             string `tfschema:"endpoint"`
	Location             string `tfschema:"location"`
	Name                 string `tfschema:"name"`
}

func (r ReplicaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicas.ValidateReplicaID
}
func (r ReplicaResource) ResourceType() string {
	return "azurerm_app_configuration_replica"
}
func (r ReplicaResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: configurationstores.ValidateConfigurationStoreID,
		},
		"location": commonschema.Location(),
		"name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ReplicaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r ReplicaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppConfiguration.ReplicasClient

			var config ReplicaResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(config.ConfigurationStoreId)
			if err != nil {
				return fmt.Errorf("parsing configuration store id %s: %+v", config.ConfigurationStoreId, err)
			}

			id := replicas.NewReplicaID(configurationStoreId.SubscriptionId, configurationStoreId.ResourceGroupName, configurationStoreId.ConfigurationStoreName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := replicas.Replica{
				Location: pointer.To(config.Location),
				Name:     pointer.To(config.Name),
			}

			// concurrent creation of replicas under one configuration store will fail
			locks.ByName(id.ConfigurationStoreName, r.ResourceType())
			defer locks.UnlockByName(id.ConfigurationStoreName, r.ResourceType())

			if err = client.CreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ReplicaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppConfiguration.ReplicasClient

			id, err := replicas.ParseReplicaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected nil model for %s", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("unexpected nil properties for %s", *id)
			}
			if resp.Model.Location == nil {
				return fmt.Errorf("unexpected nil location for %s", *id)
			}

			schema := ReplicaResourceSchema{
				ConfigurationStoreId: configurationstores.NewConfigurationStoreID(id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName).ID(),
				Endpoint:             pointer.From(resp.Model.Properties.Endpoint),
				Location:             location.Normalize(*resp.Model.Location),
				Name:                 id.ReplicaName,
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ReplicaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppConfiguration.ReplicasClient

			id, err := replicas.ParseReplicaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ConfigurationStoreName, r.ResourceType())
			defer locks.UnlockByName(id.ConfigurationStoreName, r.ResourceType())

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
