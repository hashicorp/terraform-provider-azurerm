package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobexecutions"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type MsSqlJobExecuteAction struct {
	sdk.ActionMetadata
}

var _ sdk.Action = &MsSqlJobExecuteAction{}

func newMssqlJobExecuteAction() action.Action {
	return &MsSqlJobExecuteAction{}
}

type MsSqlJobExecuteActionModel struct {
	JobID types.String `tfsdk:"job_id"`
}

func (m *MsSqlJobExecuteAction) Schema(_ context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"job_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the job to execute.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: jobexecutions.ValidateJobID,
					},
				},
			},
		},
	}
}

func (m *MsSqlJobExecuteAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_mssql_job_execute" // TODO better name?
}

func (m *MsSqlJobExecuteAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := m.Client.MSSQL.JobExecutionsClient

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	model := MsSqlJobExecuteActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	id, err := jobexecutions.ParseJobID(model.JobID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "parsing ID", err)
		return
	}

	if _, err := client.Create(ctx, *id); err != nil {
		sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Errorf("executing job: %w", err))
	}

	// TODO: should we track the actual job's execution status? This is just a fire and forget at the moment
	// tracking would require manually polling/sending get reqs on the execution ID returned by the create call

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("triggered job execution on %s", id),
	})
}

func (m *MsSqlJobExecuteAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	m.Defaults(ctx, request, response)
}
