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

type BuildServiceListBuildsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Build
}

type BuildServiceListBuildsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Build
}

type BuildServiceListBuildsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceListBuildsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceListBuilds ...
func (c AppPlatformClient) BuildServiceListBuilds(ctx context.Context, id BuildServiceId) (result BuildServiceListBuildsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceListBuildsCustomPager{},
		Path:       fmt.Sprintf("%s/builds", id.ID()),
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
		Values *[]Build `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceListBuildsComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceListBuildsComplete(ctx context.Context, id BuildServiceId) (BuildServiceListBuildsCompleteResult, error) {
	return c.BuildServiceListBuildsCompleteMatchingPredicate(ctx, id, BuildOperationPredicate{})
}

// BuildServiceListBuildsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceListBuildsCompleteMatchingPredicate(ctx context.Context, id BuildServiceId, predicate BuildOperationPredicate) (result BuildServiceListBuildsCompleteResult, err error) {
	items := make([]Build, 0)

	resp, err := c.BuildServiceListBuilds(ctx, id)
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

	result = BuildServiceListBuildsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
