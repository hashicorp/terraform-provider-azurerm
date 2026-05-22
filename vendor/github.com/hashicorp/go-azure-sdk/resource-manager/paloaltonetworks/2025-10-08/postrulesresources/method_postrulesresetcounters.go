package postrulesresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostRulesresetCountersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RuleCounterReset
}

type PostRulesresetCountersOperationOptions struct {
	FirewallName *string
}

func DefaultPostRulesresetCountersOperationOptions() PostRulesresetCountersOperationOptions {
	return PostRulesresetCountersOperationOptions{}
}

func (o PostRulesresetCountersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PostRulesresetCountersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PostRulesresetCountersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FirewallName != nil {
		out.Append("firewallName", fmt.Sprintf("%v", *o.FirewallName))
	}
	return &out
}

// PostRulesresetCounters ...
func (c PostRulesResourcesClient) PostRulesresetCounters(ctx context.Context, id PostRuleId, options PostRulesresetCountersOperationOptions) (result PostRulesresetCountersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/resetCounters", id.ID()),
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

	var model RuleCounterReset
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
