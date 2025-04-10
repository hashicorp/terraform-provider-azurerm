package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/guestagents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel struct {
	ScopedResourceId   string       `tfschema:"scoped_resource_id"`
	Credential         []Credential `tfschema:"credential"`
	ProvisioningAction string       `tfschema:"provisioning_action"`
}

type Credential struct {
	Username string `tfschema:"username"`
	Password string `tfschema:"password"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource struct{}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent"
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scoped_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"credential": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"password": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"provisioning_action": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(guestagents.ProvisioningActionInstall),
			ValidateFunc: validation.StringInSlice(guestagents.PossibleValuesForProvisioningAction(), false),
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.GuestAgents

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(model.ScopedResourceId)

			existing, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := guestagents.GuestAgent{
				Properties: &guestagents.GuestAgentProperties{
					Credentials:        expandSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(model.Credential),
					ProvisioningAction: pointer.To(guestagents.ProvisioningAction(model.ProvisioningAction)),
				},
			}

			if err := client.CreateThenPoll(ctx, commonids.NewScopeID(id.Scope), parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.GuestAgents

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel{
				ScopedResourceId: id.Scope,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Credential = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(props.Credentials, metadata.ResourceData.Get("credential.0.password").(string))
					state.ProvisioningAction = string(pointer.From(props.ProvisioningAction))
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.GuestAgents

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, commonids.NewScopeID(id.Scope)); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(input []Credential) *guestagents.GuestCredential {
	if len(input) == 0 {
		return nil
	}

	credential := input[0]

	return &guestagents.GuestCredential{
		Username: credential.Username,
		Password: credential.Password,
	}
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(input *guestagents.GuestCredential, password string) []Credential {
	result := make([]Credential, 0)
	if input == nil {
		return result
	}

	return append(result, Credential{
		Username: input.Username,
		Password: password,
	})
}
