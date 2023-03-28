package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SAPVirtualInstanceBasedDiscoveryConfigModel struct {
	Name                     string                       `tfschema:"name"`
	ResourceGroupName        string                       `tfschema:"resource_group_name"`
	Location                 string                       `tfschema:"location"`
	Configuration            []Configuration              `tfschema:"configuration"`
	Environment              string                       `tfschema:"environment"`
	Identity                 []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName string                       `tfschema:"managed_resource_group_name"`
	SapProduct               string                       `tfschema:"sap_product"`
	Tags                     map[string]string            `tfschema:"tags"`
}

type Configuration struct {
	CentralServerVmId         string `tfschema:"central_server_vm_id"`
	ManagedStorageAccountName string `tfschema:"managed_storage_account_name"`
}

type SAPVirtualInstanceBasedDiscoveryConfigResource struct{}

var _ sdk.ResourceWithUpdate = SAPVirtualInstanceBasedDiscoveryConfigResource{}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) ResourceType() string {
	return "azurerm_sap_virtual_instance_based_discovery_config"
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) ModelObject() interface{} {
	return &SAPVirtualInstanceBasedDiscoveryConfigModel{}
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"central_server_vm_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"managed_storage_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"environment": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPEnvironmentTypeNonProd),
				string(sapvirtualinstances.SAPEnvironmentTypeProd),
			}, false),
		},

		"sap_product": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPProductTypeECC),
				string(sapvirtualinstances.SAPProductTypeOther),
				string(sapvirtualinstances.SAPProductTypeSFourHANA),
			}, false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SAPVirtualInstanceBasedDiscoveryConfigModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPVirtualInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapvirtualinstances.NewSapVirtualInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Properties: sapvirtualinstances.SAPVirtualInstanceProperties{
					Configuration: expandConfiguration(model.Configuration),
					Environment:   sapvirtualinstances.SAPEnvironmentType(model.Environment),
					SapProduct:    sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

			if v := model.ManagedResourceGroupName; v != "" {
				parameters.Properties.ManagedResourceGroupConfiguration = &sapvirtualinstances.ManagedRGConfiguration{
					Name: utils.String(v),
				}
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SAPVirtualInstanceBasedDiscoveryConfigModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			parameters := &sapvirtualinstances.UpdateSAPVirtualInstanceRequest{}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if _, err := client.Update(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := SAPVirtualInstanceBasedDiscoveryConfigModel{
				Name:              id.SapVirtualInstanceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identity, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			properties := &model.Properties

			discoveryConfiguration := properties.Configuration.(sapvirtualinstances.DiscoveryConfiguration)
			state.Configuration = flattenConfiguration(&discoveryConfiguration)
			state.Environment = string(properties.Environment)
			state.SapProduct = string(properties.SapProduct)

			if v := properties.ManagedResourceGroupConfiguration; v != nil && v.Name != nil {
				state.ManagedResourceGroupName = *v.Name
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandConfiguration(input []Configuration) *sapvirtualinstances.DiscoveryConfiguration {
	if len(input) == 0 {
		return nil
	}

	configuration := &input[0]

	result := sapvirtualinstances.DiscoveryConfiguration{}

	if v := configuration.CentralServerVmId; v != "" {
		result.CentralServerVMId = utils.String(v)
	}

	if v := configuration.ManagedStorageAccountName; v != "" {
		result.ManagedRgStorageAccountName = utils.String(v)
	}

	return &result
}

func flattenConfiguration(input *sapvirtualinstances.DiscoveryConfiguration) []Configuration {
	if input == nil {
		return nil
	}

	result := Configuration{}

	if v := input.CentralServerVMId; v != nil {
		result.CentralServerVmId = *v
	}

	if v := input.ManagedRgStorageAccountName; v != nil {
		result.ManagedStorageAccountName = *v
	}

	return []Configuration{
		result,
	}
}
