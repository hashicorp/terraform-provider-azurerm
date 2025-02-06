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

type BuildServiceAgentPoolListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BuildServiceAgentPoolResource
}

type BuildServiceAgentPoolListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BuildServiceAgentPoolResource
}

type BuildServiceAgentPoolListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceAgentPoolListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceAgentPoolList ...
func (c AppPlatformClient) BuildServiceAgentPoolList(ctx context.Context, id BuildServiceId) (result BuildServiceAgentPoolListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceAgentPoolListCustomPager{},
		Path:       fmt.Sprintf("%s/agentPools", id.ID()),
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
		Values *[]BuildServiceAgentPoolResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceAgentPoolListComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceAgentPoolListComplete(ctx context.Context, id BuildServiceId) (BuildServiceAgentPoolListCompleteResult, error) {
	return c.BuildServiceAgentPoolListCompleteMatchingPredicate(ctx, id, BuildServiceAgentPoolResourceOperationPredicate{})
}

// BuildServiceAgentPoolListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceAgentPoolListCompleteMatchingPredicate(ctx context.Context, id BuildServiceId, predicate BuildServiceAgentPoolResourceOperationPredicate) (result BuildServiceAgentPoolListCompleteResult, err error) {
	items := make([]BuildServiceAgentPoolResource, 0)

	resp, err := c.BuildServiceAgentPoolList(ctx, id)
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

	result = BuildServiceAgentPoolListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
