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

type BuildServiceListSupportedStacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SupportedStackResource
}

type BuildServiceListSupportedStacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SupportedStackResource
}

type BuildServiceListSupportedStacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceListSupportedStacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceListSupportedStacks ...
func (c AppPlatformClient) BuildServiceListSupportedStacks(ctx context.Context, id BuildServiceId) (result BuildServiceListSupportedStacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceListSupportedStacksCustomPager{},
		Path:       fmt.Sprintf("%s/supportedStacks", id.ID()),
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
		Values *[]SupportedStackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceListSupportedStacksComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceListSupportedStacksComplete(ctx context.Context, id BuildServiceId) (BuildServiceListSupportedStacksCompleteResult, error) {
	return c.BuildServiceListSupportedStacksCompleteMatchingPredicate(ctx, id, SupportedStackResourceOperationPredicate{})
}

// BuildServiceListSupportedStacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceListSupportedStacksCompleteMatchingPredicate(ctx context.Context, id BuildServiceId, predicate SupportedStackResourceOperationPredicate) (result BuildServiceListSupportedStacksCompleteResult, err error) {
	items := make([]SupportedStackResource, 0)

	resp, err := c.BuildServiceListSupportedStacks(ctx, id)
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

	result = BuildServiceListSupportedStacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
