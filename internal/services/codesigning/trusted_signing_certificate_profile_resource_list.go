// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package codesigning

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	TrustedSigningCertificateProfileListResource struct{}
	TrustedSigningCertificateProfileListModel    struct {
		TrustedSigningAccountId types.String `tfsdk:"trusted_signing_account_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(TrustedSigningCertificateProfileListResource)

func (TrustedSigningCertificateProfileListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = TrustedSigningCertificateProfileResource{}.ResourceType()
}

func (TrustedSigningCertificateProfileListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(TrustedSigningCertificateProfileResource{})
}

func (TrustedSigningCertificateProfileListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"trusted_signing_account_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: certificateprofiles.ValidateCodeSigningAccountID},
				},
			},
		},
	}
}

func (TrustedSigningCertificateProfileListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.CodeSigning.Client.CertificateProfiles

	var data TrustedSigningCertificateProfileListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	accountId, err := certificateprofiles.ParseCodeSigningAccountID(data.TrustedSigningAccountId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Trusted Signing Account ID for `%s`", TrustedSigningCertificateProfileResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListByCodeSigningAccountComplete(ctx, *accountId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", TrustedSigningCertificateProfileResource{}.ResourceType()), err)
		return
	}

	r := TrustedSigningCertificateProfileResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := certificateprofiles.ParseCertificateProfileIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Certificate Profile ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := flattenTrustedSigningCertificateProfileResource(rmd, id, &item); err != nil {
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
