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

type BuildServiceBuilderListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BuilderResource
}

type BuildServiceBuilderListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BuilderResource
}

type BuildServiceBuilderListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceBuilderListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceBuilderList ...
func (c AppPlatformClient) BuildServiceBuilderList(ctx context.Context, id BuildServiceId) (result BuildServiceBuilderListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceBuilderListCustomPager{},
		Path:       fmt.Sprintf("%s/builders", id.ID()),
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
		Values *[]BuilderResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceBuilderListComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceBuilderListComplete(ctx context.Context, id BuildServiceId) (BuildServiceBuilderListCompleteResult, error) {
	return c.BuildServiceBuilderListCompleteMatchingPredicate(ctx, id, BuilderResourceOperationPredicate{})
}

// BuildServiceBuilderListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceBuilderListCompleteMatchingPredicate(ctx context.Context, id BuildServiceId, predicate BuilderResourceOperationPredicate) (result BuildServiceBuilderListCompleteResult, err error) {
	items := make([]BuilderResource, 0)

	resp, err := c.BuildServiceBuilderList(ctx, id)
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

	result = BuildServiceBuilderListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
