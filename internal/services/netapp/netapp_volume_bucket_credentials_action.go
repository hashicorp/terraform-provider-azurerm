// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type NetAppVolumeBucketCredentialsAction struct {
	sdk.ActionMetadata
}

var _ sdk.Action = &NetAppVolumeBucketCredentialsAction{}

func newNetAppVolumeBucketCredentialsAction() action.Action {
	return &NetAppVolumeBucketCredentialsAction{}
}

type NetAppVolumeBucketCredentialsActionModel struct {
	BucketID          types.String `tfsdk:"bucket_id"`
	KeyPairExpiryDays types.Int64  `tfsdk:"key_pair_expiry_days"`
	Timeout           types.String `tfsdk:"timeout"`
}

func (a *NetAppVolumeBucketCredentialsAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_netapp_volume_bucket_credentials"
}

func (a *NetAppVolumeBucketCredentialsAction) Schema(_ context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Generates an S3 access key / secret key pair for an Azure NetApp Files Volume Bucket and stores it in the Azure Key Vault configured on the parent bucket.",
		MarkdownDescription: "Generates an S3 access key / secret key pair for an Azure NetApp Files Volume Bucket and stores it in the Azure Key Vault configured on the parent bucket.",
		Attributes: map[string]schema.Attribute{
			"bucket_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ARM ID of the NetApp Volume Bucket the credentials apply to.",
				MarkdownDescription: "The ARM ID of the NetApp Volume Bucket the credentials apply to.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: buckets.ValidateBucketID,
					},
				},
			},

			"key_pair_expiry_days": schema.Int64Attribute{
				Required:            true,
				Description:         "Number of days the generated key pair is valid for. Must be at least `1`.",
				MarkdownDescription: "Number of days the generated key pair is valid for. Must be at least `1`.",
			},

			"timeout": schema.StringAttribute{
				Optional:            true,
				Description:         "Timeout duration for the bucket credentials generation to complete. Defaults to `1h`.",
				MarkdownDescription: "Timeout duration for the bucket credentials generation to complete. Defaults to `1h`.",
			},
		},
	}
}

func (a *NetAppVolumeBucketCredentialsAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := a.Client.NetApp.BucketsClient

	model := NetAppVolumeBucketCredentialsActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	timeout := 60 * time.Minute
	if t := model.Timeout; !t.IsNull() {
		duration, err := time.ParseDuration(t.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(response, "parsing `timeout`", err)
			return
		}
		timeout = duration
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id, err := buckets.ParseBucketID(model.BucketID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("parsing bucket id %s", model.BucketID.ValueString()), err)
		return
	}

	expiryDays := model.KeyPairExpiryDays.ValueInt64()
	if expiryDays < 1 {
		sdk.SetResponseErrorDiagnostic(response, "invalid `key_pair_expiry_days`", fmt.Errorf("`key_pair_expiry_days` must be at least 1, got %d", expiryDays))
		return
	}

	input := buckets.BucketCredentialsExpiry{
		KeyPairExpiryDays: pointer.To(expiryDays),
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("generating Key Vault-backed credentials for %s", id),
	})

	if err := client.GenerateAkvCredentialsThenPoll(ctx, pointer.From(id), input); err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("generating credentials in key vault for %s", id), err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("credentials successfully generated and stored in Key Vault for %s", id),
	})
}

func (a *NetAppVolumeBucketCredentialsAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)
}
