// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementLockListResource struct{}

type ManagementLockListModel struct {
	Scope             types.String `tfsdk:"scope"`
	SubscriptionId    types.String `tfsdk:"subscription_id"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	ResourceId        types.String `tfsdk:"resource_id"`
}

var _ sdk.FrameworkListWrappedResource = &ManagementLockListResource{}

func (ManagementLockListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_management_lock"
}

func (ManagementLockListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceManagementLock()
}

func (ManagementLockListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the Subscription to query locks for (supports either raw UUID or /subscriptions/{subscriptionId}). Queries locks applied directly at the subscription level.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.Any(validation.IsUUID, commonids.ValidateSubscriptionID),
					},
					stringvalidator.ConflictsWith(
						path.MatchRoot("resource_group_name"),
						path.MatchRoot("resource_id"),
						path.MatchRoot("scope"),
					),
				},
			},
			"resource_group_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the Resource Group to query locks for. Queries locks applied directly at the resource group level.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
					stringvalidator.ConflictsWith(
						path.MatchRoot("subscription_id"),
						path.MatchRoot("resource_id"),
						path.MatchRoot("scope"),
					),
				},
			},
			"resource_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the Resource to query locks for. Queries locks applied directly at the resource level.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateScopeID,
					},
					stringvalidator.ConflictsWith(
						path.MatchRoot("subscription_id"),
						path.MatchRoot("resource_group_name"),
						path.MatchRoot("scope"),
					),
				},
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Description: "The Scope to query locks for (e.g. a Management Group or arbitrary resource scope).",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateScopeID,
					},
					stringvalidator.ConflictsWith(
						path.MatchRoot("subscription_id"),
						path.MatchRoot("resource_group_name"),
						path.MatchRoot("resource_id"),
					),
				},
			},
		},
	}
}

func (ManagementLockListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Resource.LocksClient

	var data ManagementLockListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var items []managementlocks.ManagementLockObject

	switch {
	case !data.Scope.IsNull():
		scope := data.Scope.ValueString()
		scopeId := commonids.NewScopeID(scope)
		resp, err := client.ListByScopeComplete(ctx, scopeId, managementlocks.ListByScopeOperationOptions{})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing management locks in scope `%s`", scope), err)
			return
		}
		items = resp.Items

	case !data.ResourceId.IsNull():
		resourceId := data.ResourceId.ValueString()
		resScopeId := commonids.NewScopeID(resourceId)
		resp, err := client.ListAtResourceLevelComplete(ctx, resScopeId, managementlocks.ListAtResourceLevelOperationOptions{})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing management locks at resource `%s`", resourceId), err)
			return
		}
		items = resp.Items

	case !data.ResourceGroupName.IsNull():
		rgName := data.ResourceGroupName.ValueString()
		rgId := commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, rgName)
		resp, err := client.ListAtResourceGroupLevelComplete(ctx, rgId, managementlocks.ListAtResourceGroupLevelOperationOptions{})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing management locks in resource group `%s`", rgName), err)
			return
		}
		items = resp.Items

	default:
		subInput := metadata.Client.Account.SubscriptionId
		if !data.SubscriptionId.IsNull() {
			subInput = data.SubscriptionId.ValueString()
			if parsed, err := commonids.ParseSubscriptionID(subInput); err == nil {
				subInput = parsed.SubscriptionId
			}
		}
		subId := commonids.NewSubscriptionID(subInput)
		resp, err := client.ListAtSubscriptionLevelComplete(ctx, subId, managementlocks.ListAtSubscriptionLevelOperationOptions{})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing management locks at subscription `%s`", subId), err)
			return
		}
		items = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := managementlocks.ParseScopedLockID(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Management Lock ID", err)
				return
			}

			rd := resourceManagementLock().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceManagementLockFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting resource data for `azurerm_management_lock`", err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
