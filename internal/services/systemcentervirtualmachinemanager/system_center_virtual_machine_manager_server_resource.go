package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerServerModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	CustomLocationId  string            `tfschema:"custom_location_id"`
	Fqdn              string            `tfschema:"fqdn"`
	Credential        []Credential      `tfschema:"credential"`
	Port              int               `tfschema:"port"`
	Tags              map[string]string `tfschema:"tags"`
}

type Credential struct {
	Username string `tfschema:"username"`
	Password string `tfschema:"password"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerServerResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerServerResource{}

type SystemCenterVirtualMachineManagerServerResource struct{}

func (r SystemCenterVirtualMachineManagerServerResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerServerModel{}
}

func (r SystemCenterVirtualMachineManagerServerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return vmmservers.ValidateVMmServerID
}

func (r SystemCenterVirtualMachineManagerServerResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_server"
}

func (r SystemCenterVirtualMachineManagerServerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemCenterVirtualMachineManagerServerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"tags": commonschema.Tags(),
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerServerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			var model SystemCenterVirtualMachineManagerServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := vmmservers.NewVMmServerID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &vmmservers.VMMServer{
				Location: location.Normalize(model.Location),
				ExtendedLocation: vmmservers.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				},
				Properties: vmmservers.VMMServerProperties{
					Credentials: expandSystemCenterVirtualMachineManagerServerCredential(model.Credential),
					Fqdn:        model.Fqdn,
				},
				Tags: pointer.To(model.Tags),
			}

			if model.Port != 0 {
				parameters.Properties.Port = utils.Int64(int64(model.Port))
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			id, err := vmmservers.ParseVMmServerID(metadata.ResourceData.Id())
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

			state := SystemCenterVirtualMachineManagerServerModel{}
			if model := resp.Model; model != nil {
				state.Name = id.VmmServerName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.Fqdn = model.Properties.Fqdn
				state.Credential = flattenSystemCenterVirtualMachineManagerServerCredential(model.Properties.Credentials)
				state.Tags = pointer.From(model.Tags)

				if v := model.Properties.Port; v != nil {
					state.Port = int(*v)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			id, err := vmmservers.ParseVMmServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := existing.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("custom_location_id") {
				parameters.ExtendedLocation = vmmservers.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				}
			}

			if metadata.ResourceData.HasChange("fqdn") {
				parameters.Properties.Fqdn = model.Fqdn
			}

			if metadata.ResourceData.HasChange("port") {
				parameters.Properties.Port = utils.Int64(int64(model.Port))
			}

			if metadata.ResourceData.HasChange("credential") {
				parameters.Properties.Credentials = expandSystemCenterVirtualMachineManagerServerCredential(model.Credential)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			id, err := vmmservers.ParseVMmServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, vmmservers.DeleteOperationOptions{Force: pointer.To(vmmservers.ForceTrue)}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandSystemCenterVirtualMachineManagerServerCredential(input []Credential) *vmmservers.VMMCredential {
	if len(input) == 0 {
		return nil
	}

	credential := &input[0]

	result := &vmmservers.VMMCredential{
		Username: pointer.To(credential.Username),
		Password: pointer.To(credential.Password),
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerServerCredential(input *vmmservers.VMMCredential) []Credential {
	result := make([]Credential, 0)
	if input == nil {
		return result
	}

	credential := Credential{
		Username: pointer.From(input.Username),
		Password: pointer.From(input.Password),
	}

	return append(result, credential)
}
