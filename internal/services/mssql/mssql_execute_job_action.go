package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobexecutions"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/custompollers"
)

type MsSqlExecuteJobAction struct {
	sdk.ActionMetadata
}

var _ sdk.Action = &MsSqlExecuteJobAction{}

func newMssqlJobExecuteAction() action.Action {
	return &MsSqlExecuteJobAction{}
}

type MsSqlJobExecuteActionModel struct {
	JobID             types.String `tfsdk:"job_id"`
	WaitForCompletion types.Bool   `tfsdk:"wait_for_completion"`
	Timeout           types.String `tfsdk:"timeout"`
}

func (m *MsSqlExecuteJobAction) Schema(_ context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"job_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the job to execute.",
				MarkdownDescription: "The ID of the job to execute.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: jobexecutions.ValidateJobID,
					},
				},
			},

			"wait_for_completion": schema.BoolAttribute{
				Optional:            true,
				Description:         "Whether to poll the job execution for completion. Defaults to `false`",
				MarkdownDescription: "Whether to poll the job execution for completion. Defaults to `false`",
			},

			"timeout": schema.StringAttribute{
				Optional:            true,
				Description:         "Timeout duration for the action to complete. Defaults to `15m`.",
				MarkdownDescription: "Timeout duration for the action to complete. Defaults to `15m`.",
			},
		},
	}
}

func (m *MsSqlExecuteJobAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_mssql_execute_job"
}

func (m *MsSqlExecuteJobAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := m.Client.MSSQL.JobExecutionsClient
	jobStepExecutionsClient := m.Client.MSSQL.JobStepExecutionsClient

	model := MsSqlJobExecuteActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	timeout := 15 * time.Minute
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

	id, err := jobexecutions.ParseJobID(model.JobID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "parsing ID", err)
		return
	}

	resp, err := client.Create(ctx, *id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Errorf("executing job: %w", err))
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("triggered execution of `%s`", id.ID()),
	})

	if model.WaitForCompletion.ValueBool() {
		response.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("waiting for completion on `%s`", id.ID()),
		})

		if resp.Model == nil || resp.Model.Id == nil {
			sdk.SetResponseErrorDiagnostic(response, "waiting for completion", "unable to retrieve execution ID")
			return
		}

		executionId, err := jobexecutions.ParseExecutionID(*resp.Model.Id)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(response, "waiting for completion", err)
			return
		}

		executionPoller := custompollers.NewMssqlJobExecutionStatusPoller(client, jobStepExecutionsClient, *executionId)
		poller := pollers.NewPoller(executionPoller, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "waiting for completion", err)
			return
		}

		response.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("execution of `%s` completed", id.ID()),
		})
	}
}

func (m *MsSqlExecuteJobAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	m.Defaults(ctx, request, response)
}
