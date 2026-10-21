// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managedidentity

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	FederatedIdentityCredentialListResource struct{}
	FederatedIdentityCredentialListModel    struct {
		UserAssignedIdentityId types.String `tfsdk:"user_assigned_identity_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(FederatedIdentityCredentialListResource)

func (FederatedIdentityCredentialListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = FederatedIdentityCredentialResource{}.ResourceType()
}

func (FederatedIdentityCredentialListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(FederatedIdentityCredentialResource{})
}

func (FederatedIdentityCredentialListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_assigned_identity_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateUserAssignedIdentityID},
				},
			},
		},
	}
}

func (FederatedIdentityCredentialListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

	var data FederatedIdentityCredentialListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	userAssignedIdentityId, err := commonids.ParseUserAssignedIdentityID(data.UserAssignedIdentityId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing User Assigned Identity ID for `%s`", FederatedIdentityCredentialResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, *userAssignedIdentityId, federatedidentitycredentials.DefaultListOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", FederatedIdentityCredentialResource{}.ResourceType()), err)
		return
	}

	r := FederatedIdentityCredentialResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Federated Identity Credential ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
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
