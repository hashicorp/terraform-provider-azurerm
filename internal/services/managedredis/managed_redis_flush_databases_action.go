package managedredis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type ManagedRedisFlushDatabasesAction struct {
	sdk.ActionMetadata
}

// var _ sdk.Action = &ManagedRedisFlushDatabasesAction{}
func newManagedRedisFlushDatabasesAction() action.Action {
	return &ManagedRedisFlushDatabasesAction{}
}

type ManagedRedisFlushDatabasesActionActionModel struct {
	ManagedRedisDatabaseId types.String   `tfsdk:"managed_redis_database_id"`
	LinkedDatabaseIds      []types.String `tfsdk:"linked_database_ids"`
	Timeout                types.String   `tfsdk:"timeout"`
}

func (m *ManagedRedisFlushDatabasesAction) Schema(_ context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"managed_redis_database_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the Managed Redis Database where the keys will be flushed.",
				MarkdownDescription: "The ID of the Managed Redis Database where the keys will be flushed.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: redisenterprise.ValidateRedisEnterpriseID,
					},
				},
			},

			"linked_database_ids": schema.ListAttribute{
				Optional:            true,
				Description:         "The IDs of any Linked Databases where the keys will be flushed.",
				MarkdownDescription: "The IDs of any Linked Databases where the keys will be flushed.",
				ElementType:         types.StringType,
				Validators: []validator.List{
					typehelpers.WrappedListValidator{
						Func: redisenterprise.ValidateRedisEnterpriseID,
					},
				},
			},

			"timeout": schema.StringAttribute{
				Optional:            true,
				Description:         "Timeout duration for the Managed Redis Databases Flush action to complete. Defaults to 30m.",
				MarkdownDescription: "Timeout duration for the Managed Redis Databases Flush action to complete. Defaults to 30m.",
			},
		},
	}
}

func (m *ManagedRedisFlushDatabasesAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_managed_redis_databases_flush"
}

func (m *ManagedRedisFlushDatabasesAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := m.Client.ManagedRedis.Client

	model := ManagedRedisFlushDatabasesActionActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	timeout := 30 * time.Minute
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

	id, err := redisenterprise.ParseDatabaseID(model.ManagedRedisDatabaseId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "id parsing error", err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("flushing %s", id),
	})

	linkedDatabaseIds := make([]string, 0)
	for _, linkedDatabaseId := range model.LinkedDatabaseIds {
		linkedDatabaseIds = append(linkedDatabaseIds, linkedDatabaseId.ValueString())
	}

	payload := redisenterprise.FlushParameters{
		Ids: pointer.To(linkedDatabaseIds),
	}

	if err = client.DatabasesFlushThenPoll(ctx, *id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("flushing database for %s", id), err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("flushing completed for %s", id.ID()),
	})
}

func (m *ManagedRedisFlushDatabasesAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	m.Defaults(ctx, request, response)
}
