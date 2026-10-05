package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildServiceListBuildResultsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BuildResult
}

type BuildServiceListBuildResultsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BuildResult
}

type BuildServiceListBuildResultsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceListBuildResultsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceListBuildResults ...
func (c AppPlatformClient) BuildServiceListBuildResults(ctx context.Context, id BuildId) (result BuildServiceListBuildResultsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceListBuildResultsCustomPager{},
		Path:       fmt.Sprintf("%s/results", id.ID()),
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
		Values *[]BuildResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceListBuildResultsComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceListBuildResultsComplete(ctx context.Context, id BuildId) (BuildServiceListBuildResultsCompleteResult, error) {
	return c.BuildServiceListBuildResultsCompleteMatchingPredicate(ctx, id, BuildResultOperationPredicate{})
}

// BuildServiceListBuildResultsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceListBuildResultsCompleteMatchingPredicate(ctx context.Context, id BuildId, predicate BuildResultOperationPredicate) (result BuildServiceListBuildResultsCompleteResult, err error) {
	items := make([]BuildResult, 0)

	resp, err := c.BuildServiceListBuildResults(ctx, id)
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

	result = BuildServiceListBuildResultsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
