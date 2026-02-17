package redis

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedisFirewallRuleListResource struct{}

type RedisFirewallRuleListModel struct {
	RedisCacheId types.String `tfsdk:"redis_cache_id"`
}

var _ sdk.FrameworkListWrappedResource = new(RedisFirewallRuleListResource)

func (RedisFirewallRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceRedisFirewallRule()
}

func (r RedisFirewallRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = redisFirewallRuleResourceName
}

func (r RedisFirewallRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"redis_cache_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: redisfirewallrules.ValidateRediID,
					},
				},
			},
		},
	}
}

func (RedisFirewallRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {

	client := metadata.Client.Redis.FirewallRulesClient

	// Read the list config data into the model
	var data RedisFirewallRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	// Initialize a list for the results of the API request
	results := make([]redisfirewallrules.RedisFirewallRule, 0)

	cacheId, err := redisfirewallrules.ParseRediIDInsensitively(data.RedisCacheId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Redis Cache ID for `%s`", redisFirewallRuleResourceName), err)
		return
	}

	resp, err := client.FirewallRulesListComplete(ctx, pointer.From(cacheId))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", redisFirewallRuleResourceName), err)
		return
	}
	results = resp.Items

	// Define the function that will push results into the stream
	stream.Results = func(push func(list.ListResult) bool) {
		for _, rules := range results {

			// Initialize a new result object for each resource in the list
			result := request.NewListResult(ctx)

			// Set the display name of the item as the resource name
			result.DisplayName = pointer.From(rules.Name)

			// Create a new ResourceData object to hold the state of the resource
			rd := resourceRedisFirewallRule().Data(&terraform.InstanceState{})

			// Set the ID of the resource for the ResourceData object
			// API is returning /Redis/ with capital "R", so need to parse insensitive
			id, err := redisfirewallrules.ParseFirewallRuleIDInsensitively(pointer.From(rules.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Redis Firewall Rule ID", err)
				return
			}
			rd.SetId(id.ID())

			// Use the resource flatten function to set the attributes into the resource state
			if err := resourceRedisFirewallRuleFlatten(rd, id, &rules); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", redisFirewallRuleResourceName), err)
				return
			}

			// Convert and set the identity and resource state into the result
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
