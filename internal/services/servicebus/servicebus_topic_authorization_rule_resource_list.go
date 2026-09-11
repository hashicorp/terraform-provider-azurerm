package servicebus

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/topicsauthorizationrule"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceBusTopicAuthorizationRuleListResource struct{}

type ServiceBusTopicAuthorizationRuleListModel struct {
	TopicId types.String `tfsdk:"servicebus_topic_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ServiceBusTopicAuthorizationRuleListResource)

func (ServiceBusTopicAuthorizationRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = serviceBusTopicAuthorizationRuleResourceName
}

func (ServiceBusTopicAuthorizationRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceServiceBusTopicAuthorizationRule()
}

func (ServiceBusTopicAuthorizationRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"servicebus_topic_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: topicsauthorizationrule.ValidateTopicID},
				},
			},
		},
	}
}

func (ServiceBusTopicAuthorizationRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ServiceBus.TopicsAuthClient

	var data ServiceBusTopicAuthorizationRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := topicsauthorizationrule.ParseTopicID(data.TopicId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Topic ID for `%s`", serviceBusTopicAuthorizationRuleResourceName), err)
		return
	}

	resp, err := client.TopicsListAuthorizationRulesComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusTopicAuthorizationRuleResourceName), err)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceServiceBusTopicAuthorizationRule().Data(&terraform.InstanceState{})
			id, err := topicsauthorizationrule.ParseTopicAuthorizationRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", serviceBusTopicAuthorizationRuleResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceServiceBusTopicAuthorizationRuleFlatten(ctx, client, rd, id, &item, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", serviceBusTopicAuthorizationRuleResourceName), err)
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
