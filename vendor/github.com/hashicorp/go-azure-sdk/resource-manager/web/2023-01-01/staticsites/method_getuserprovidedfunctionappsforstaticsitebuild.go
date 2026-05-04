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

type GetUserProvidedFunctionAppsForStaticSiteBuildOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteUserProvidedFunctionAppARMResource
}

type GetUserProvidedFunctionAppsForStaticSiteBuildCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteUserProvidedFunctionAppARMResource
}

type GetUserProvidedFunctionAppsForStaticSiteBuildCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetUserProvidedFunctionAppsForStaticSiteBuildCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetUserProvidedFunctionAppsForStaticSiteBuild ...
func (c StaticSitesClient) GetUserProvidedFunctionAppsForStaticSiteBuild(ctx context.Context, id BuildId) (result GetUserProvidedFunctionAppsForStaticSiteBuildOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetUserProvidedFunctionAppsForStaticSiteBuildCustomPager{},
		Path:       fmt.Sprintf("%s/userProvidedFunctionApps", id.ID()),
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
		Values *[]StaticSiteUserProvidedFunctionAppARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetUserProvidedFunctionAppsForStaticSiteBuildComplete retrieves all the results into a single object
func (c StaticSitesClient) GetUserProvidedFunctionAppsForStaticSiteBuildComplete(ctx context.Context, id BuildId) (GetUserProvidedFunctionAppsForStaticSiteBuildCompleteResult, error) {
	return c.GetUserProvidedFunctionAppsForStaticSiteBuildCompleteMatchingPredicate(ctx, id, StaticSiteUserProvidedFunctionAppARMResourceOperationPredicate{})
}

// GetUserProvidedFunctionAppsForStaticSiteBuildCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetUserProvidedFunctionAppsForStaticSiteBuildCompleteMatchingPredicate(ctx context.Context, id BuildId, predicate StaticSiteUserProvidedFunctionAppARMResourceOperationPredicate) (result GetUserProvidedFunctionAppsForStaticSiteBuildCompleteResult, err error) {
	items := make([]StaticSiteUserProvidedFunctionAppARMResource, 0)

	resp, err := c.GetUserProvidedFunctionAppsForStaticSiteBuild(ctx, id)
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

	result = GetUserProvidedFunctionAppsForStaticSiteBuildCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
