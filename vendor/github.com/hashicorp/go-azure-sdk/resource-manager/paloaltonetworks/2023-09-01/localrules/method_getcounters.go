package localrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetCountersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RuleCounter
}

type GetCountersOperationOptions struct {
	FirewallName *string
}

func DefaultGetCountersOperationOptions() GetCountersOperationOptions {
	return GetCountersOperationOptions{}
}

func (o GetCountersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCountersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GetCountersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FirewallName != nil {
		out.Append("firewallName", fmt.Sprintf("%v", *o.FirewallName))
	}
	return &out
}

// GetCounters ...
func (c LocalRulesClient) GetCounters(ctx context.Context, id LocalRuleId, options GetCountersOperationOptions) (result GetCountersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/getCounters", id.ID()),
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
