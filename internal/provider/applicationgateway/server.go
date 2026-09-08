// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package applicationgateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const resourceName = "azurerm_application_gateway"

// NewServerFactory preserves unchanged Application Gateway blocks after the SDK
// reconstructs the plan. CustomizeDiff runs before that reconstruction.
func NewServerFactory(provider *schema.Provider) func() tfprotov5.ProviderServer {
	return func() tfprotov5.ProviderServer {
		return wrapServer(provider.GRPCProvider(), provider.ResourcesMap[resourceName])
	}
}

func wrapServer(server tfprotov5.ProviderServer, resource *schema.Resource) tfprotov5.ProviderServer {
	if resource == nil {
		return server
	}
	return &planServer{
		ProviderServer: server,
		resourceType:   resource.CoreConfigSchema().ImpliedType(),
		resource:       resource,
	}
}

type planServer struct {
	tfprotov5.ProviderServer
	resourceType cty.Type
	resource     *schema.Resource
}

func (s *planServer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	resp, err := s.ProviderServer.PlanResourceChange(ctx, req)
	if err != nil || !canPreserveBlocks(req, resp) {
		return resp, err
	}
	prior, err := decode(req.PriorState, s.resourceType)
	if err != nil {
		return nil, fmt.Errorf("decoding Application Gateway prior state: %w", err)
	}
	planned, err := decode(resp.PlannedState, s.resourceType)
	if err != nil {
		return nil, fmt.Errorf("decoding Application Gateway planned state: %w", err)
	}
	if !isExistingGatewayUpdate(prior, planned) {
		return resp, nil
	}

	preserved := preserveUnchangedBlocks(prior, planned, s.resource)
	if preserved.RawEquals(planned) {
		return resp, nil
	}
	encoded, err := msgpack.Marshal(preserved, s.resourceType)
	if err != nil {
		return nil, fmt.Errorf("encoding Application Gateway planned state: %w", err)
	}
	resp.PlannedState = &tfprotov5.DynamicValue{MsgPack: encoded}
	return resp, nil
}

func canPreserveBlocks(req *tfprotov5.PlanResourceChangeRequest, resp *tfprotov5.PlanResourceChangeResponse) bool {
	if req.TypeName != resourceName || req.PriorState == nil || resp == nil || resp.PlannedState == nil {
		return false
	}
	if len(resp.RequiresReplace) > 0 || resp.Deferred != nil {
		return false
	}
	for _, diagnostic := range resp.Diagnostics {
		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			return false
		}
	}
	return true
}

func isExistingGatewayUpdate(prior, planned cty.Value) bool {
	if prior.IsNull() || planned.IsNull() || !prior.IsKnown() || !planned.IsKnown() {
		return false
	}
	id := prior.GetAttr("id")
	return id.IsKnown() && id.RawEquals(planned.GetAttr("id"))
}

func decode(value *tfprotov5.DynamicValue, typ cty.Type) (cty.Value, error) {
	if len(value.MsgPack) > 0 {
		return msgpack.Unmarshal(value.MsgPack, typ)
	}
	return ctyjson.Unmarshal(value.JSON, typ)
}
