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

type ListStaticSiteBuildFunctionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteFunctionOverviewARMResource
}

type ListStaticSiteBuildFunctionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteFunctionOverviewARMResource
}

type ListStaticSiteBuildFunctionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListStaticSiteBuildFunctionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListStaticSiteBuildFunctions ...
func (c StaticSitesClient) ListStaticSiteBuildFunctions(ctx context.Context, id BuildId) (result ListStaticSiteBuildFunctionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListStaticSiteBuildFunctionsCustomPager{},
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

// ListStaticSiteBuildFunctionsComplete retrieves all the results into a single object
func (c StaticSitesClient) ListStaticSiteBuildFunctionsComplete(ctx context.Context, id BuildId) (ListStaticSiteBuildFunctionsCompleteResult, error) {
	return c.ListStaticSiteBuildFunctionsCompleteMatchingPredicate(ctx, id, StaticSiteFunctionOverviewARMResourceOperationPredicate{})
}

// ListStaticSiteBuildFunctionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) ListStaticSiteBuildFunctionsCompleteMatchingPredicate(ctx context.Context, id BuildId, predicate StaticSiteFunctionOverviewARMResourceOperationPredicate) (result ListStaticSiteBuildFunctionsCompleteResult, err error) {
	items := make([]StaticSiteFunctionOverviewARMResource, 0)

	resp, err := c.ListStaticSiteBuildFunctions(ctx, id)
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

	result = ListStaticSiteBuildFunctionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
