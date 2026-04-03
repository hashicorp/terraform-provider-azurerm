package managednetwork

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SettingsRuleListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OutboundRuleBasicResource
}

type SettingsRuleListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OutboundRuleBasicResource
}

type SettingsRuleListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SettingsRuleListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SettingsRuleList ...
func (c ManagedNetworkClient) SettingsRuleList(ctx context.Context, id WorkspaceId) (result SettingsRuleListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SettingsRuleListCustomPager{},
		Path:       fmt.Sprintf("%s/outboundRules", id.ID()),
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
		Values *[]OutboundRuleBasicResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SettingsRuleListComplete retrieves all the results into a single object
func (c ManagedNetworkClient) SettingsRuleListComplete(ctx context.Context, id WorkspaceId) (SettingsRuleListCompleteResult, error) {
	return c.SettingsRuleListCompleteMatchingPredicate(ctx, id, OutboundRuleBasicResourceOperationPredicate{})
}

// SettingsRuleListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedNetworkClient) SettingsRuleListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate OutboundRuleBasicResourceOperationPredicate) (result SettingsRuleListCompleteResult, err error) {
	items := make([]OutboundRuleBasicResource, 0)

	resp, err := c.SettingsRuleList(ctx, id)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = SettingsRuleListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
