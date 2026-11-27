// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	BillingInvoiceSectionListResource struct{}
	BillingInvoiceSectionListModel    struct {
		BillingProfileId types.String `tfsdk:"billing_profile_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(BillingInvoiceSectionListResource)

func (BillingInvoiceSectionListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = BillingInvoiceSectionResource{}.ResourceType()
}

func (BillingInvoiceSectionListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(BillingInvoiceSectionResource{})
}

func (BillingInvoiceSectionListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	// invoice sections can only be listed within a billing profile
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"billing_profile_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: invoicesection.ValidateBillingProfileID},
				},
			},
		},
	}
}

func (BillingInvoiceSectionListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Billing.InvoiceSectionClient

	r := BillingInvoiceSectionResource{}

	var data BillingInvoiceSectionListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	billingProfileId, err := invoicesection.ParseBillingProfileID(data.BillingProfileId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Billing Profile ID for `%s`", r.ResourceType()), err)
		return
	}

	resp, err := client.ListByBillingProfileComplete(ctx, *billingProfileId, invoicesection.DefaultListByBillingProfileOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := invoicesection.ParseInvoiceSectionIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Invoice Section ID", err)
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
