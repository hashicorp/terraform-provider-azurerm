package dnsforwardingrulesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByVirtualNetworkOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualNetworkDnsForwardingRuleset
}

type ListByVirtualNetworkCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualNetworkDnsForwardingRuleset
}

type ListByVirtualNetworkOperationOptions struct {
	Top *int64
}

func DefaultListByVirtualNetworkOperationOptions() ListByVirtualNetworkOperationOptions {
	return ListByVirtualNetworkOperationOptions{}
}

func (o ListByVirtualNetworkOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByVirtualNetworkOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByVirtualNetworkOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByVirtualNetwork ...
func (c DnsForwardingRulesetsClient) ListByVirtualNetwork(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions) (result ListByVirtualNetworkOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/listDnsForwardingRulesets", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]VirtualNetworkDnsForwardingRuleset `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByVirtualNetworkComplete retrieves all the results into a single object
func (c DnsForwardingRulesetsClient) ListByVirtualNetworkComplete(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions) (ListByVirtualNetworkCompleteResult, error) {
	return c.ListByVirtualNetworkCompleteMatchingPredicate(ctx, id, options, VirtualNetworkDnsForwardingRulesetOperationPredicate{})
}

// ListByVirtualNetworkCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DnsForwardingRulesetsClient) ListByVirtualNetworkCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions, predicate VirtualNetworkDnsForwardingRulesetOperationPredicate) (result ListByVirtualNetworkCompleteResult, err error) {
	items := make([]VirtualNetworkDnsForwardingRuleset, 0)

	resp, err := c.ListByVirtualNetwork(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListByVirtualNetworkCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
