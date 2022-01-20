package containers

import (
	"context"
	"fmt"
	"time"

	legacyacr "github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2019-06-01-preview/containerregistry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	networkparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkvalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryAgentPoolResource struct{}

var _ sdk.ResourceWithUpdate = ContainerRegistryAgentPoolResource{}

type AgentPoolModel struct {
	Name                string                 `tfschema:"name"`
	ContainerRegistryId string                 `tfschema:"container_registry_id"`
	Tier                string                 `tfschema:"tier"`
	AgentCount          int                    `tfschema:"agent_count"`
	OS                  string                 `tfschema:"os"`
	SubnetId            string                 `tfschema:"subnet_id"`
	Tags                map[string]interface{} `tfschema:"tags"`
}

func (r ContainerRegistryAgentPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AgentPoolName,
		},
		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RegistryID,
		},
		"tier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"S1",
				"S2",
				"S3",
				"I6",
			}, false),
		},
		"agent_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 100),
		},
		"os": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(legacyacr.Linux),
			ValidateFunc: validation.StringInSlice([]string{string(legacyacr.Linux)}, false),
		},
		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: networkvalidate.SubnetID,
		},
		"tags": tags.ForceNewSchema(),
	}
}

func (r ContainerRegistryAgentPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerRegistryAgentPoolResource) ResourceType() string {
	return "azurerm_container_registry_agent_pool"
}

func (r ContainerRegistryAgentPoolResource) ModelObject() interface{} {
	return &AgentPoolModel{}
}

func (r ContainerRegistryAgentPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AgentPoolID
}

func (r ContainerRegistryAgentPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			acrClient := metadata.Client.Containers.RegistriesClient
			client := metadata.Client.Containers.ContainerRegistryAgentPoolsClient

			var model AgentPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			rid, err := parse.RegistryID(model.ContainerRegistryId)
			if err != nil {
				return err
			}
			existingAcr, err := acrClient.Get(ctx, rid.ResourceGroup, rid.Name)
			if err != nil {
				return fmt.Errorf("getting %s: %v", rid, err)
			}

			id := parse.NewAgentPoolID(rid.SubscriptionId, rid.ResourceGroup, rid.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := legacyacr.AgentPool{
				AgentPoolProperties: &legacyacr.AgentPoolProperties{
					Count: utils.Int32(int32(model.AgentCount)),
					Tier:  &model.Tier,
					Os:    legacyacr.OS(model.OS),
				},
				Location: utils.String(location.NormalizeNilable(existingAcr.Location)),
				Tags:     tags.Expand(model.Tags),
			}

			if model.SubnetId != "" {
				subnetId, _ := networkparse.SubnetID(model.SubnetId)
				params.AgentPoolProperties.VirtualNetworkSubnetResourceID = utils.String(subnetId.ID())
			}

			future, err := client.Create(ctx, id.ResourceGroup, id.RegistryName, id.Name, params)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContainerRegistryAgentPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.AgentPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AgentPoolModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Containers.ContainerRegistryAgentPoolsClient

			patch := legacyacr.AgentPoolUpdateParameters{}
			if metadata.ResourceData.HasChange("tags") {
				patch.Tags = tags.Expand(state.Tags)
			}
			if metadata.ResourceData.HasChange("agent_count") {
				patch.AgentPoolPropertiesUpdateParameters = &legacyacr.AgentPoolPropertiesUpdateParameters{Count: utils.Int32(int32(state.AgentCount))}
			}

			future, err := client.Update(ctx, id.ResourceGroup, id.RegistryName, id.Name, patch)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ContainerRegistryAgentPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryAgentPoolsClient
			id, err := parse.AgentPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			rid := parse.NewRegistryID(id.SubscriptionId, id.ResourceGroup, id.RegistryName)
			model := AgentPoolModel{
				Name:                id.Name,
				ContainerRegistryId: rid.ID(),
				Tags:                tags.Flatten(existing.Tags),
			}

			if props := existing.AgentPoolProperties; props != nil {
				if tier := props.Tier; tier != nil {
					model.Tier = *tier
				}
				if count := props.Count; count != nil {
					model.AgentCount = int(*count)
				}
				model.OS = string(props.Os)
				if subnetIdPtr := props.VirtualNetworkSubnetResourceID; subnetIdPtr != nil {
					subnetId, err := networkparse.SubnetID(*subnetIdPtr)
					if err != nil {
						return err
					}
					model.SubnetId = subnetId.ID()
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ContainerRegistryAgentPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryAgentPoolsClient

			id, err := parse.AgentPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for removal of %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
