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

type PreRulesrefreshCountersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type PreRulesrefreshCountersOperationOptions struct {
	FirewallName *string
}

func DefaultPreRulesrefreshCountersOperationOptions() PreRulesrefreshCountersOperationOptions {
	return PreRulesrefreshCountersOperationOptions{}
}

func (o PreRulesrefreshCountersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PreRulesrefreshCountersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PreRulesrefreshCountersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FirewallName != nil {
		out.Append("firewallName", fmt.Sprintf("%v", *o.FirewallName))
	}
	return &out
}

// PreRulesrefreshCounters ...
func (c PreRulesResourcesClient) PreRulesrefreshCounters(ctx context.Context, id PreRuleId, options PreRulesrefreshCountersOperationOptions) (result PreRulesrefreshCountersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/refreshCounters", id.ID()),
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

	return
}
