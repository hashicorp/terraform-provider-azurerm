package raiblocklists

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RaiBlocklistItemsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RaiBlocklistItem
}

type RaiBlocklistItemsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RaiBlocklistItem
}

type RaiBlocklistItemsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RaiBlocklistItemsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RaiBlocklistItemsList ...
func (c RaiBlocklistsClient) RaiBlocklistItemsList(ctx context.Context, id RaiBlocklistId) (result RaiBlocklistItemsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RaiBlocklistItemsListCustomPager{},
		Path:       fmt.Sprintf("%s/raiBlocklistItems", id.ID()),
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
		Values *[]RaiBlocklistItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RaiBlocklistItemsListComplete retrieves all the results into a single object
func (c RaiBlocklistsClient) RaiBlocklistItemsListComplete(ctx context.Context, id RaiBlocklistId) (RaiBlocklistItemsListCompleteResult, error) {
	return c.RaiBlocklistItemsListCompleteMatchingPredicate(ctx, id, RaiBlocklistItemOperationPredicate{})
}

// RaiBlocklistItemsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RaiBlocklistsClient) RaiBlocklistItemsListCompleteMatchingPredicate(ctx context.Context, id RaiBlocklistId, predicate RaiBlocklistItemOperationPredicate) (result RaiBlocklistItemsListCompleteResult, err error) {
	items := make([]RaiBlocklistItem, 0)

	resp, err := c.RaiBlocklistItemsList(ctx, id)
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

	result = RaiBlocklistItemsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
