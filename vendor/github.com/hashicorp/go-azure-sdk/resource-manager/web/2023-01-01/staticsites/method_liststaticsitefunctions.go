package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListStaticSiteFunctionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteFunctionOverviewARMResource
}

type ListStaticSiteFunctionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteFunctionOverviewARMResource
}

type ListStaticSiteFunctionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListStaticSiteFunctionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListStaticSiteFunctions ...
func (c StaticSitesClient) ListStaticSiteFunctions(ctx context.Context, id StaticSiteId) (result ListStaticSiteFunctionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListStaticSiteFunctionsCustomPager{},
		Path:       fmt.Sprintf("%s/functions", id.ID()),
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
		Values *[]StaticSiteFunctionOverviewARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListStaticSiteFunctionsComplete retrieves all the results into a single object
func (c StaticSitesClient) ListStaticSiteFunctionsComplete(ctx context.Context, id StaticSiteId) (ListStaticSiteFunctionsCompleteResult, error) {
	return c.ListStaticSiteFunctionsCompleteMatchingPredicate(ctx, id, StaticSiteFunctionOverviewARMResourceOperationPredicate{})
}

// ListStaticSiteFunctionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) ListStaticSiteFunctionsCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate StaticSiteFunctionOverviewARMResourceOperationPredicate) (result ListStaticSiteFunctionsCompleteResult, err error) {
	items := make([]StaticSiteFunctionOverviewARMResource, 0)

	resp, err := c.ListStaticSiteFunctions(ctx, id)
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

	result = ListStaticSiteFunctionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
