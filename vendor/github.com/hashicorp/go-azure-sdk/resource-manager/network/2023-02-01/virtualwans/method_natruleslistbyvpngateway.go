package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NatRulesListByVpnGatewayOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnGatewayNatRule
}

type NatRulesListByVpnGatewayCompleteResult struct {
	Items []VpnGatewayNatRule
}

// NatRulesListByVpnGateway ...
func (c VirtualWANsClient) NatRulesListByVpnGateway(ctx context.Context, id VpnGatewayId) (result NatRulesListByVpnGatewayOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/natRules", id.ID()),
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
		Values *[]VpnGatewayNatRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NatRulesListByVpnGatewayComplete retrieves all the results into a single object
func (c VirtualWANsClient) NatRulesListByVpnGatewayComplete(ctx context.Context, id VpnGatewayId) (NatRulesListByVpnGatewayCompleteResult, error) {
	return c.NatRulesListByVpnGatewayCompleteMatchingPredicate(ctx, id, VpnGatewayNatRuleOperationPredicate{})
}

// NatRulesListByVpnGatewayCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) NatRulesListByVpnGatewayCompleteMatchingPredicate(ctx context.Context, id VpnGatewayId, predicate VpnGatewayNatRuleOperationPredicate) (result NatRulesListByVpnGatewayCompleteResult, err error) {
	items := make([]VpnGatewayNatRule, 0)

	resp, err := c.NatRulesListByVpnGateway(ctx, id)
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

	result = NatRulesListByVpnGatewayCompleteResult{
		Items: items,
	}
	return
}
