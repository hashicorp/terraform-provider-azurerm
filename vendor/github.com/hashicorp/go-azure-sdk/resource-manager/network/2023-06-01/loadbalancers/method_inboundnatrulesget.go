package loadbalancers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundNatRulesGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *InboundNatRule
}

type InboundNatRulesGetOperationOptions struct {
	Expand *string
}

func DefaultInboundNatRulesGetOperationOptions() InboundNatRulesGetOperationOptions {
	return InboundNatRulesGetOperationOptions{}
}

func (o InboundNatRulesGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o InboundNatRulesGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o InboundNatRulesGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// InboundNatRulesGet ...
func (c LoadBalancersClient) InboundNatRulesGet(ctx context.Context, id InboundNatRuleId, options InboundNatRulesGetOperationOptions) (result InboundNatRulesGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          id.ID(),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
