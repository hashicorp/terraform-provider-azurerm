package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vminstanceguestagents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel struct {
	ScopedResourceId   string       `tfschema:"scoped_resource_id"`
	Credential         []Credential `tfschema:"credential"`
	HttpsProxy         string       `tfschema:"https_proxy"`
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
			ValidateFunc: computevalidate.HybridMachineID,
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
						ForceNew: true,
					},

					"password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						ForceNew:  true,
						Sensitive: true,
					},
				},
			},
		},

		"https_proxy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"provisioning_action": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(vminstanceguestagents.PossibleValuesForProvisioningAction(), false),
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
			client := metadata.Client.SystemCenterVirtualMachineManager.VMInstanceGuestAgents

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

			parameters := vminstanceguestagents.GuestAgent{
				Properties: vminstanceguestagents.GuestAgentProperties{
					Credentials:     expandSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(model.Credential),
					HTTPProxyConfig: &vminstanceguestagents.HTTPProxyConfiguration{},
				},
			}

			if v := model.HttpsProxy; v != "" {
				parameters.Properties.HTTPProxyConfig.HTTPSProxy = utils.String(v)
			}

			if v := model.ProvisioningAction; v != "" {
				parameters.Properties.ProvisioningAction = pointer.To(vminstanceguestagents.ProvisioningAction(model.ProvisioningAction))
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
			client := metadata.Client.SystemCenterVirtualMachineManager.VMInstanceGuestAgents

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

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel{}
			state.ScopedResourceId = id.Scope
			if model := resp.Model; model != nil {
				state.Credential = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(model.Properties.Credentials)
				state.ProvisioningAction = string(pointer.From(model.Properties.ProvisioningAction))

				if v := model.Properties.HTTPProxyConfig; v != nil {
					state.HttpsProxy = pointer.From(v.HTTPSProxy)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMInstanceGuestAgents

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := existing.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			parameters.Properties.Credentials = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(model.Credential)
			parameters.Properties.ProvisioningAction = pointer.To(vminstanceguestagents.ProvisioningAction(model.ProvisioningAction))
			parameters.Properties.HTTPProxyConfig = &vminstanceguestagents.HTTPProxyConfiguration{
				HTTPSProxy: utils.String(model.HttpsProxy),
			}

			if err := client.CreateThenPoll(ctx, commonids.NewScopeID(id.Scope), *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMInstanceGuestAgents

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

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(input []Credential) *vminstanceguestagents.GuestCredential {
	if len(input) == 0 {
		return nil
	}

	credential := input[0]

	return &vminstanceguestagents.GuestCredential{
		Username: credential.Username,
		Password: credential.Password,
	}
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentCredential(input *vminstanceguestagents.GuestCredential) []Credential {
	result := make([]Credential, 0)
	if input == nil {
		return result
	}

	return append(result, Credential{
		Username: input.Username,
		Password: input.Password,
	})
}
