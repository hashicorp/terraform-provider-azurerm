// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobcontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobservices"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageContainerListResource struct{}

func (r StorageContainerListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceStorageContainer()
}

var _ sdk.FrameworkListWrappedResource = &StorageContainerListResource{}

type StorageContainerListModel struct {
	StorageAccountID types.String `tfsdk:"storage_account_id"`
}

func (r StorageContainerListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = StorageContainerResourceName
}

func (r StorageContainerListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"storage_account_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateStorageAccountID,
					},
				},
			},
		},
	}
}

func (r StorageContainerListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Storage.ResourceManager.BlobServices

	var data StorageContainerListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	storageAccountID, err := commonids.ParseStorageAccountID(data.StorageAccountID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Parent ID for `%s`", StorageContainerResourceName), err)
		return
	}

	resp, err := client.BlobContainersListComplete(ctx, *storageAccountID, blobservices.DefaultBlobContainersListOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", StorageContainerResourceName), err)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(ctx, deadline)
		defer cancel()

		for _, c := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(c.Name)

			id, err := commonids.ParseStorageContainerIDInsensitively(*c.Id)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Container ID", err)
				return
			}

			rd := resourceStorageContainer().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err = resourceStorageContainerFlatten(ctx, rd, id, r.listItemToBlobContainer(c), metadata.Client, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "encoding Resource data", err)
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

// listItemToBlobContainer converts a returned container response to the type expected by the resource's flatten function
func (r StorageContainerListResource) listItemToBlobContainer(item blobservices.ListContainerItem) *blobcontainers.BlobContainer {
	result := &blobcontainers.BlobContainer{
		Name: item.Name,
		Type: item.Type,
	}

	if props := item.Properties; props != nil {
		result.Properties = &blobcontainers.ContainerProperties{
			DefaultEncryptionScope:      props.DefaultEncryptionScope,
			DenyEncryptionScopeOverride: props.DenyEncryptionScopeOverride,
			HasImmutabilityPolicy:       props.HasImmutabilityPolicy,
			HasLegalHold:                props.HasLegalHold,
			Metadata:                    props.Metadata,
			PublicAccess:                pointer.ToEnum[blobcontainers.PublicAccess](pointer.FromEnum(props.PublicAccess)),
		}
	}

	return result
}
