// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	ScopedResourceId   string `tfschema:"scoped_resource_id"`
	Username           string `tfschema:"username"`
	Password           string `tfschema:"password"`
	ProvisioningAction string `tfschema:"provisioning_action"`
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
					Credentials: &guestagents.GuestCredential{
						Username: model.Username,
						Password: model.Password,
					},
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
					if v := props.Credentials; v != nil {
						state.Username = v.Username
						state.Password = metadata.ResourceData.Get("password").(string)
					}

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
