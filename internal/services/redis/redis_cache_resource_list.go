package redis

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedisCacheListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(RedisCacheListResource)

func (RedisCacheListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceRedisCache()
}

func (r RedisCacheListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = redisCacheResourceName
}

func (RedisCacheListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {

	client := metadata.Client.Redis.RedisResourcesClient

	// retrieve the deadline from the supplied context
	deadline, ok := ctx.Deadline()
	if !ok {
		// This *should* never happen given the List Wrapper instantiates a context with a timeout
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	// Read the list config data into the model
	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	// Initialize a list for the results of the API request
	results := make([]redisresources.RedisResource, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	// Make the request based on which list parameters have been set in the config
	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.RedisListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", redisCacheResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.RedisListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", redisCacheResourceName), err)
			return
		}

		results = resp.Items
	}

	// Define the function that will push results into the stream
	stream.Results = func(push func(list.ListResult) bool) {

		// Instantiate a new context based on the deadline retrieved earlier
		deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, redis := range results {

			// Initialize a new result object for each resource in the list
			result := request.NewListResult(deadlineCtx)

			// Set the display name of the item as the resource name
			result.DisplayName = pointer.From(redis.Name)

			// Create a new ResourceData object to hold the state of the resource
			rd := resourceRedisCache().Data(&terraform.InstanceState{})

			// Set the ID of the resource for the ResourceData object
			// API is returning /Redis/ with capital "R", so need to parse insensitive
			id, err := redisresources.ParseRediIDInsensitively(pointer.From(redis.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Redis Cache ID", err)
				return
			}
			rd.SetId(id.ID())

			// Use the resource flatten function to set the attributes into the resource state
			if err := resourceRedisCacheFlatten(rd, id, &redis, deadlineCtx, metadata.Client); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", redisCacheResourceName), err)
				return
			}

			// Convert and set the identity and resource state into the result
			sdk.EncodeListResult(deadlineCtx, rd, &result)
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
