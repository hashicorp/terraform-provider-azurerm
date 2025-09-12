// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ListResourceWithRawV5Schemas = &StorageAccountListResource{}

type StorageAccountListResource struct {
	sdk.ListResourceMetadata
}

type StorageAccountListModel struct {
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	SubscriptionId    types.String `tfsdk:"subscription_id"`
}

func NewStorageAccountListResource() list.ListResource {
	return &StorageAccountListResource{}
}

func (r *StorageAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = storageAccountResourceName
}

func (r *StorageAccountListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, response *list.RawV5SchemaResponse) {
	res := resourceStorageAccount()
	response.ProtoV5Schema = res.ProtoSchema(ctx)()
	response.ProtoV5IdentitySchema = res.ProtoIdentitySchema(ctx)()
}

func (r *StorageAccountListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	// Storage Accounts support listing by Resource Group Name or Subscription ID, both are optional here and we default
	// to the local subscription ID if nothing is provided
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"resource_group_name": listschema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},
			"subscription_id": listschema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.IsUUID,
					},
				},
			},
		},
	}
}

func (r *StorageAccountListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	storageClient := r.Client.Storage.ResourceManager
	client := storageClient.StorageAccounts

	ctx, cancel := context.WithTimeout(ctx, time.Minute*60) // TODO - Can/should we make this user configurable?
	defer cancel()

	var data StorageAccountListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	listResults := make([]storageaccounts.StorageAccount, 0)
	subscriptionID := r.SubscriptionId
	if data.SubscriptionId.ValueString() != "" {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case data.ResourceGroupName.ValueString() != "":
		resourceGroupId := commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString())
		resp, err := client.ListByResourceGroupComplete(ctx, resourceGroupId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", storageAccountResourceName), err)
			return
		}

		listResults = resp.Items

	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", storageAccountResourceName), err)
			return
		}

		listResults = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		// TODO - Do we need to handle limiting the results to ListRequest.Limit?
		variableTimeout := time.Duration(5*len(listResults)) * time.Minute
		ctx, cancel := context.WithTimeout(context.Background(), variableTimeout)
		defer cancel()
		for _, account := range listResults {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(account.Name)
			id, err := commonids.ParseStorageAccountID(*account.Id)
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "parsing storage account id", err)
				return
			}

			saResource := resourceStorageAccount()

			rd := saResource.Data(&terraform.InstanceState{})

			rd.SetId(id.ID())

			if err := resourceStorageAccountFlatten(ctx, rd, *id, pointer.To(account), r.Client); err != nil {
				sdk.SetResponseWarningDiagnostic(stream, "encoding resource data", err)
				// Not erroring here as best effort on additional API call(s) made by the flatten function can error out
				// when we have enough data to perform the import.
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting identity data", err)
				return
			}

			tfTypeResource, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Resource State data", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting resource data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}

func (r *StorageAccountListResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Defaults(request, response)
}
