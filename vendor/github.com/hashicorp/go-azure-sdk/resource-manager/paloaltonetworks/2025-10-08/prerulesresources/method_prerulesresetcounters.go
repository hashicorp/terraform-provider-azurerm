package prerulesresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PreRulesresetCountersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RuleCounterReset
}

type PreRulesresetCountersOperationOptions struct {
	FirewallName *string
}

func DefaultPreRulesresetCountersOperationOptions() PreRulesresetCountersOperationOptions {
	return PreRulesresetCountersOperationOptions{}
}

func (o PreRulesresetCountersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PreRulesresetCountersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PreRulesresetCountersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FirewallName != nil {
		out.Append("firewallName", fmt.Sprintf("%v", *o.FirewallName))
	}
	return &out
}

// PreRulesresetCounters ...
func (c PreRulesResourcesClient) PreRulesresetCounters(ctx context.Context, id PreRuleId, options PreRulesresetCountersOperationOptions) (result PreRulesresetCountersOperationResponse, err error) {
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
