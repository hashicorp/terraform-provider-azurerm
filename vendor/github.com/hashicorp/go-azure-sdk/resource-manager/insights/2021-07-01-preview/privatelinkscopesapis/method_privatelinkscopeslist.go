package privatelinkscopesapis

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

type PrivateLinkScopesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AzureMonitorPrivateLinkScope
}

type PrivateLinkScopesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AzureMonitorPrivateLinkScope
}

type PrivateLinkScopesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PrivateLinkScopesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PrivateLinkScopesList ...
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesList(ctx context.Context, id commonids.SubscriptionId) (result PrivateLinkScopesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PrivateLinkScopesListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/privateLinkScopes", id.ID()),
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
		Values *[]AzureMonitorPrivateLinkScope `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PrivateLinkScopesListComplete retrieves all the results into a single object
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListComplete(ctx context.Context, id commonids.SubscriptionId) (PrivateLinkScopesListCompleteResult, error) {
	return c.PrivateLinkScopesListCompleteMatchingPredicate(ctx, id, AzureMonitorPrivateLinkScopeOperationPredicate{})
}

// PrivateLinkScopesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AzureMonitorPrivateLinkScopeOperationPredicate) (result PrivateLinkScopesListCompleteResult, err error) {
	items := make([]AzureMonitorPrivateLinkScope, 0)

	resp, err := c.PrivateLinkScopesList(ctx, id)
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

	result = PrivateLinkScopesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
