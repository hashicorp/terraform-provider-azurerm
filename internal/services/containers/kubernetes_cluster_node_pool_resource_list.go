package containers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesClusterNodePoolListResource struct{}

type KubernetesClusterNodePoolListModel struct {
	KubernetesClusterId types.String `tfsdk:"kubernetes_cluster_id"`
}

var _ sdk.FrameworkListWrappedResource = new(KubernetesClusterNodePoolListResource)

func (KubernetesClusterNodePoolListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureKubernetesClusterNodePoolResourceName
}

func (KubernetesClusterNodePoolListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceKubernetesClusterNodePool()
}

func (KubernetesClusterNodePoolListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"kubernetes_cluster_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateKubernetesClusterID},
				},
			},
		},
	}
}

func (KubernetesClusterNodePoolListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Containers.AgentPoolsClient

	var data KubernetesClusterNodePoolListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := commonids.ParseKubernetesClusterID(data.KubernetesClusterId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Kubernetes Cluster ID for `%s`", azureKubernetesClusterNodePoolResourceName), err)
		return
	}

	resp, err := client.ListComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureKubernetesClusterNodePoolResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceKubernetesClusterNodePool().Data(&terraform.InstanceState{})
			id, err := agentpools.ParseAgentPoolIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureKubernetesClusterNodePoolResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceKubernetesClusterNodePoolFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureKubernetesClusterNodePoolResourceName), err)
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
