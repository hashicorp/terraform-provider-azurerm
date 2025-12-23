// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupinstanceresources"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

const (
	StopProtection   = "stop_protection"
	ResumeProtection = "resume_protection"
	SuspendBackups   = "suspend_backups"
	ResumeBackups    = "resume_backups"
)

type DataProtectionBackupInstanceProtectAction struct {
	sdk.ActionMetadata
}

var _ sdk.Action = &DataProtectionBackupInstanceProtectAction{}

func newDataProtectionBackupInstanceProtectAction() action.Action {
	return &DataProtectionBackupInstanceProtectAction{}
}

type DataProtectionBackupInstanceProtectActionModel struct {
	BackupInstanceId types.String `tfsdk:"backup_instance_id"`
	Action           types.String `tfsdk:"protect_action"`
}

func (v *DataProtectionBackupInstanceProtectAction) Schema(_ context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"backup_instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the backup instance on which to perform the action.",
				MarkdownDescription: "The ID of the backup instance on which to perform the action.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: backupinstanceresources.ValidateBackupInstanceID,
					},
				},
			},

			"protect_action": schema.StringAttribute{
				Required:            true,
				Description:         "The protect state action to take on this backup instance. Possible values include `stop_protection`,`resume_protection`, `suspend_backups`, and `resume_backups`.",
				MarkdownDescription: "The protect state action to take on this backup instance. Possible values include `stop_protection`,`resume_protection`, `suspend_backups`, and `resume_backups`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						StopProtection,
						ResumeProtection,
						SuspendBackups,
						ResumeBackups,
					),
				},
			},
		},
	}
}

func (v *DataProtectionBackupInstanceProtectAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_data_protection_backup_instance_protect"
}

func (v *DataProtectionBackupInstanceProtectAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := v.Client.DataProtection.BackupInstanceClient

	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	model := DataProtectionBackupInstanceProtectActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	id, err := backupinstanceresources.ParseBackupInstanceID(model.BackupInstanceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "parsing id", err)
		return
	}

	protectAction := model.Action.ValueString()

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("invoking %s on %s", protectAction, id.BackupInstanceName),
	})

	switch protectAction {
	case StopProtection:
		if err := client.BackupInstancesStopProtectionThenPoll(ctx, *id, backupinstanceresources.StopProtectionRequest{}, backupinstanceresources.DefaultBackupInstancesStopProtectionOperationOptions()); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("stopping protection %s: %+v", id, err))
			return
		}

	case ResumeProtection:
		if err := client.BackupInstancesResumeProtectionThenPoll(ctx, *id); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("resuming protection %s: %+v", id, err))
			return
		}

	case SuspendBackups:
		if err := client.BackupInstancesSuspendBackupsThenPoll(ctx, *id, backupinstanceresources.SuspendBackupRequest{}, backupinstanceresources.DefaultBackupInstancesSuspendBackupsOperationOptions()); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("suspending backups %s: %+v", id, err))
			return
		}

	case ResumeBackups:
		if err := client.BackupInstancesResumeBackupsThenPoll(ctx, *id); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("resuming backups %s: %+v", id, err))
			return
		}
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("action %s on %s completed", protectAction, id.BackupInstanceName),
	})
}

func (v *DataProtectionBackupInstanceProtectAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	v.Defaults(ctx, request, response)
}
