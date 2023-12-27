package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPDiscoveryVirtualInstanceModel struct {
	Name                     string                       `tfschema:"name"`
	ResourceGroupName        string                       `tfschema:"resource_group_name"`
	Location                 string                       `tfschema:"location"`
	DiscoveryConfiguration   []DiscoveryConfiguration     `tfschema:"discovery_configuration"`
	Environment              string                       `tfschema:"environment"`
	Identity                 []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName string                       `tfschema:"managed_resource_group_name"`
	SapProduct               string                       `tfschema:"sap_product"`
	Tags                     map[string]string            `tfschema:"tags"`
}

type DiscoveryConfiguration struct {
	CentralServerVmId         string `tfschema:"central_server_virtual_machine_id"`
	ManagedStorageAccountName string `tfschema:"managed_storage_account_name"`
}

type WorkloadsSAPDiscoveryVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPDiscoveryVirtualInstanceResource{}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_discovery_virtual_instance"
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPDiscoveryVirtualInstanceModel{}
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SAPVirtualInstanceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"discovery_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"central_server_virtual_machine_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateVirtualMachineID,
					},

					"managed_storage_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: storageValidate.StorageAccountName,
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

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPDiscoveryVirtualInstanceModel
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

			identity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Properties: sapvirtualinstances.SAPVirtualInstanceProperties{
					Environment: sapvirtualinstances.SAPEnvironmentType(model.Environment),
					SapProduct:  sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

			if v := model.DiscoveryConfiguration; v != nil {
				parameters.Properties.Configuration = expandDiscoveryConfiguration(v)
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

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPDiscoveryVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
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

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Read() sdk.ResourceFunc {
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

			state := WorkloadsSAPDiscoveryVirtualInstanceModel{}
			if model := resp.Model; model != nil {
				state.Name = id.SapVirtualInstanceName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)

				identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(identity)

				props := &model.Properties
				state.Environment = string(props.Environment)
				state.SapProduct = string(props.SapProduct)
				state.Tags = pointer.From(model.Tags)

				if config := props.Configuration; config != nil {
					if v, ok := config.(sapvirtualinstances.DiscoveryConfiguration); ok {
						state.DiscoveryConfiguration = flattenDiscoveryConfiguration(v)
					}
				}

				if v := props.ManagedResourceGroupConfiguration; v != nil {
					state.ManagedResourceGroupName = pointer.From(v.Name)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
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

func expandDiscoveryConfiguration(input []DiscoveryConfiguration) *sapvirtualinstances.DiscoveryConfiguration {
	configuration := input[0]

	result := &sapvirtualinstances.DiscoveryConfiguration{
		CentralServerVMId: utils.String(configuration.CentralServerVmId),
	}

	if v := configuration.ManagedStorageAccountName; v != "" {
		result.ManagedRgStorageAccountName = utils.String(v)
	}

	return result
}

func flattenDiscoveryConfiguration(input sapvirtualinstances.DiscoveryConfiguration) []DiscoveryConfiguration {
	result := make([]DiscoveryConfiguration, 0)

	return append(result, DiscoveryConfiguration{
		CentralServerVmId:         pointer.From(input.CentralServerVMId),
		ManagedStorageAccountName: pointer.From(input.ManagedRgStorageAccountName),
	})
}
