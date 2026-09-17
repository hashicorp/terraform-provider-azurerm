// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection"
)

type ResourceGuardUnlockDeleteAction struct {
	sdk.ActionMetadata
}

var _ sdk.Action = &ResourceGuardUnlockDeleteAction{}

func newResourceGuardUnlockDeleteAction() action.Action {
	return &ResourceGuardUnlockDeleteAction{}
}

type ResourceGuardUnlockDeleteActionModel struct {
	ProtectedItemId types.String `tfsdk:"protected_item_id"`
}

func (a *ResourceGuardUnlockDeleteAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"protected_item_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the Recovery Services protected item to unlock for deletion.",
				MarkdownDescription: "The ID of the Recovery Services protected item to unlock for deletion.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: protecteditems.ValidateProtectedItemID,
					},
				},
			},
		},
	}
}

func (a *ResourceGuardUnlockDeleteAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "azurerm_resource_guard_unlock_delete"
}

func (a *ResourceGuardUnlockDeleteAction) Invoke(ctx context.Context, request action.InvokeRequest, resp *action.InvokeResponse) {
	model := ResourceGuardUnlockDeleteActionModel{}
	resp.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(model.ProtectedItemId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing protected item id", err)
		return
	}

	client := a.Client.RecoveryServices
	existing, err := client.ProtectedItemsClient.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("%s no longer exists; no unlock is required", id),
			})
			return
		}
		sdk.SetResponseErrorDiagnostic(resp, "retrieving protected item", fmt.Errorf("retrieving %s: %+v", id, err))
		return
	}

	vaultId := resourceguardproxies.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)
	proxies, err := client.ResourceGuardProxiesClient.GetComplete(ctx, vaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "retrieving Resource Guard proxies", fmt.Errorf("retrieving proxies for %s: %+v", vaultId, err))
		return
	}

	unlocked := false
	for _, proxy := range proxies.Items {
		var operationRequests []string
		if proxy.Properties != nil && proxy.Properties.ResourceGuardOperationDetails != nil {
			for _, detail := range *proxy.Properties.ResourceGuardOperationDetails {
				if pointer.From(detail.VaultCriticalOperation) != dataprotection.GuardOperationDeleteProtectedItem {
					continue
				}
				if pointer.From(detail.DefaultResourceRequest) == "" {
					sdk.SetResponseErrorDiagnostic(resp, "reading Resource Guard operation", fmt.Sprintf("the delete operation for %s has no default resource request", vaultId))
					return
				}
				operationRequests = append(operationRequests, *detail.DefaultResourceRequest)
			}
		}
		if len(operationRequests) == 0 {
			continue
		}

		proxyId, err := resourceguardproxy.ParseBackupResourceGuardProxyIDInsensitively(pointer.From(proxy.Id))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "parsing Resource Guard proxy id", err)
			return
		}

		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("unlocking deletion of %s using %s", id, proxyId),
		})
		unlock := resourceguardproxy.UnlockDeleteRequest{
			ResourceGuardOperationRequests: pointer.To(operationRequests),
			ResourceToBeDeleted:            pointer.To(id.ID()),
		}
		result, err := client.ResourceGuardProxyClient.UnlockDelete(ctx, *proxyId, unlock)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "unlocking deletion", fmt.Errorf("unlocking %s using %s: %+v", id, proxyId, err))
			return
		}
		unlocked = true
		if result.Model != nil && result.Model.UnlockDeleteExpiryTime != nil {
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("deletion of %s is unlocked until %s", id, *result.Model.UnlockDeleteExpiryTime),
			})
		}
	}

	message := fmt.Sprintf("unlocked deletion of %s", id)
	if !unlocked {
		message = fmt.Sprintf("deletion of %s is not protected by Resource Guard; no unlock is required", id)
	}
	resp.SendProgress(action.InvokeProgressEvent{Message: message})
}

func (a *ResourceGuardUnlockDeleteAction) Configure(ctx context.Context, request action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.Defaults(ctx, request, resp)
}
