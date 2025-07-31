package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type ListResource interface {
	list.ListResourceWithConfigure
}

type ListResourceMetadata struct {
	Client *clients.Client

	SubscriptionId string

	Features features.UserFeatures
}

func (r *ListResourceMetadata) Defaults(req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	//resp.Diagnostics.AddWarning("test warning", "this is a test warning")

	r.Client = c
	r.SubscriptionId = c.Account.SubscriptionId
	r.Features = c.Features
}

//// DecodeList performs a Get on the OpenRequest config and attempts to load it into the interface cfg. cfg *must* be a pointer to the struct.
//// returns true if successful, false if there is an error diagnostic raised. Any error diags are written directly to the response
//func (r *ListResourceMetadata) DecodeList(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream, cfg interface{}) bool {
//	//resp.Diagnostics.Append(req.Config.Get(ctx, cfg)...)
//	//
//	//return !resp.Diagnostics.HasError()
//}
