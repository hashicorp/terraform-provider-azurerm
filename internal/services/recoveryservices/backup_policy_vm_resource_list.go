// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/backuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	BackupProtectionPolicyVMListResource struct{}
	BackupProtectionPolicyVMListModel    struct {
		RecoveryVaultId types.String `tfsdk:"recovery_vault_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(BackupProtectionPolicyVMListResource)

func (r BackupProtectionPolicyVMListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceBackupProtectionPolicyVM()
}

func (r BackupProtectionPolicyVMListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_backup_policy_vm"
}

func (r BackupProtectionPolicyVMListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"recovery_vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: vaults.ValidateVaultID,
					},
				},
			},
		},
	}
}

func (r BackupProtectionPolicyVMListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	backupPoliciesClient := metadata.Client.RecoveryServices.BackupPoliciesClient
	protectionPoliciesClient := metadata.Client.RecoveryServices.ProtectionPoliciesClient

	var data BackupProtectionPolicyVMListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vaultId, err := vaults.ParseVaultID(data.RecoveryVaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing `recovery_vault_id` for `%s`", "azurerm_backup_policy_vm"), err)
		return
	}

	listVaultId := backuppolicies.NewVaultID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.VaultName)

	// Backup Policies are listed at the vault level and include all policy types, so we filter to Azure VM policies only.
	filter := "backupManagementType eq 'AzureIaasVM'"
	resp, err := backupPoliciesClient.ListComplete(ctx, listVaultId, backuppolicies.ListOperationOptions{Filter: &filter})
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_backup_policy_vm"), err)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", fmt.Errorf("context had no deadline"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, policy := range resp.Items {
			// skip any policies that are not Azure VM backup policies
			if _, ok := policy.Properties.(backuppolicies.AzureIaaSVMProtectionPolicy); !ok {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(policy.Name)

			id, err := protectionpolicies.ParseBackupPolicyID(pointer.From(policy.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Backup Policy ID", err)
				return
			}

			rd := resourceBackupProtectionPolicyVM().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			var model *protectionpolicies.ProtectionPolicyResource
			if request.IncludeResource {
				read, err := protectionPoliciesClient.Get(ctx, *id)
				if err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("retrieving %s", id), err)
					return
				}
				model = read.Model
			}

			if err := resourceBackupProtectionPolicyVMFlatten(rd, id, model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_backup_policy_vm"), err)
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
