// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                  = DesktopVirtualizationHostPoolRegistrationInfoResource{}
	_ sdk.ResourceWithUpdate        = DesktopVirtualizationHostPoolRegistrationInfoResource{}
	_ sdk.ResourceWithCustomizeDiff = DesktopVirtualizationHostPoolRegistrationInfoResource{}
)

type DesktopVirtualizationHostPoolRegistrationInfoResource struct{}

func (DesktopVirtualizationHostPoolRegistrationInfoResource) ModelObject() interface{} {
	return &DesktopVirtualizationHostPoolRegistrationInfoResourceModel{}
}

func (DesktopVirtualizationHostPoolRegistrationInfoResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.HostPoolRegistrationInfoID
}

func (DesktopVirtualizationHostPoolRegistrationInfoResource) ResourceType() string {
	return "azurerm_virtual_desktop_host_pool_registration_info"
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			if rd.HasChange("expiration_date") {
				if err := rd.SetNewComputed("token"); err != nil {
					return err
				}
				return nil
			}
			return nil
		},
	}
}

type DesktopVirtualizationHostPoolRegistrationInfoResourceModel struct {
	HostpoolId     string `tfschema:"hostpool_id"`
	ExpirationDate string `tfschema:"expiration_date"`
	Token          string `tfschema:"token"`
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostpool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hostpool.ValidateHostPoolID,
		},

		"expiration_date": {
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
			Required:     true,
		},
	}
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"token": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient

			var model DesktopVirtualizationHostPoolRegistrationInfoResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			hostPoolId, err := hostpool.ParseHostPoolID(model.HostpoolId)
			if err != nil {
				return err
			}

			locks.ByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())
			defer locks.UnlockByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())

			// This is a virtual resource so the last segment is hardcoded
			id := parse.NewHostPoolRegistrationInfoID(hostPoolId.SubscriptionId, hostPoolId.ResourceGroupName, hostPoolId.HostPoolName, "default")

			existing, err := client.Get(ctx, *hostPoolId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s could not be found: %s", hostPoolId, err)
				}
				return fmt.Errorf("reading %s: %s", hostPoolId, err)
			}

			tokenOperation := hostpool.RegistrationTokenOperationUpdate
			payload := hostpool.HostPoolPatch{
				Properties: &hostpool.HostPoolPatchProperties{
					RegistrationInfo: &hostpool.RegistrationInfoPatch{
						ExpirationTime:             pointer.To(model.ExpirationDate),
						RegistrationTokenOperation: &tokenOperation,
					},
				},
			}
			if _, err := client.Update(ctx, *hostPoolId, payload); err != nil {
				return fmt.Errorf("updating registration token for %s: %+v", hostPoolId, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient

			state := DesktopVirtualizationHostPoolRegistrationInfoResourceModel{}

			id, err := parse.HostPoolRegistrationInfoID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			hostPoolId := hostpool.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName)
			resp, err := client.RetrieveRegistrationToken(ctx, hostPoolId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] Registration Token was not found for %s - removing from state!", hostPoolId)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving Registration Token for %s: %+v", hostPoolId, err)
			}

			if resp.Model == nil || resp.Model.ExpirationTime == nil || resp.Model.Token == nil {
				log.Printf("HostPool is missing registration info - marking as gone")
				return metadata.MarkAsGone(id)
			}

			state.HostpoolId = hostPoolId.ID()
			state.ExpirationDate = pointer.From(resp.Model.ExpirationTime)
			state.Token = pointer.From(resp.Model.Token)

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient

			var model DesktopVirtualizationHostPoolRegistrationInfoResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			hostPoolId, err := hostpool.ParseHostPoolID(model.HostpoolId)
			if err != nil {
				return err
			}

			locks.ByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())
			defer locks.UnlockByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())

			existing, err := client.Get(ctx, *hostPoolId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s could not be found: %s", hostPoolId, err)
				}
				return fmt.Errorf("reading %s: %s", hostPoolId, err)
			}

			tokenOperation := hostpool.RegistrationTokenOperationUpdate
			payload := hostpool.HostPoolPatch{
				Properties: &hostpool.HostPoolPatchProperties{
					RegistrationInfo: &hostpool.RegistrationInfoPatch{
						ExpirationTime:             pointer.To(model.ExpirationDate),
						RegistrationTokenOperation: &tokenOperation,
					},
				},
			}
			if _, err := client.Update(ctx, *hostPoolId, payload); err != nil {
				return fmt.Errorf("updating registration token for %s: %+v", hostPoolId, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationHostPoolRegistrationInfoResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient
			id, err := parse.HostPoolRegistrationInfoID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			hostPoolId := hostpool.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName)

			locks.ByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())
			defer locks.UnlockByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())

			resp, err := client.Get(ctx, hostPoolId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", hostPoolId)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", hostPoolId, err)
			}

			regInfo, err := client.RetrieveRegistrationToken(ctx, hostPoolId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] Virtual Desktop Host Pool %q Registration Info was not found in Resource Group %q - removing from state!", id.HostPoolName, id.ResourceGroup)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving Registration Token for %s: %+v", hostPoolId, err)
			}
			if regInfo.Model == nil || regInfo.Model.ExpirationTime == nil {
				log.Printf("[INFO] RegistrationInfo for %s was nil, registrationInfo already deleted - removing from state", hostPoolId)
				return nil
			}

			tokenOperation := hostpool.RegistrationTokenOperationDelete
			payload := hostpool.HostPoolPatch{
				Properties: &hostpool.HostPoolPatchProperties{
					RegistrationInfo: &hostpool.RegistrationInfoPatch{
						RegistrationTokenOperation: &tokenOperation,
					},
				},
			}

			if _, err := client.Update(ctx, hostPoolId, payload); err != nil {
				return fmt.Errorf("removing Registration Token from %s: %+v", hostPoolId, err)
			}

			return nil
		},
	}
}
