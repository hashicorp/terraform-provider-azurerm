package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsagesListByLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type UsagesListByLocationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type UsagesListByLocationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *UsagesListByLocationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// UsagesListByLocation ...
func (c OpenapisClient) UsagesListByLocation(ctx context.Context, id LocationId) (result UsagesListByLocationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &UsagesListByLocationCustomPager{},
		Path:       fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]Usage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UsagesListByLocationComplete retrieves all the results into a single object
func (c OpenapisClient) UsagesListByLocationComplete(ctx context.Context, id LocationId) (UsagesListByLocationCompleteResult, error) {
	return c.UsagesListByLocationCompleteMatchingPredicate(ctx, id, UsageOperationPredicate{})
}

// UsagesListByLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) UsagesListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate UsageOperationPredicate) (result UsagesListByLocationCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.UsagesListByLocation(ctx, id)
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

	result = UsagesListByLocationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
