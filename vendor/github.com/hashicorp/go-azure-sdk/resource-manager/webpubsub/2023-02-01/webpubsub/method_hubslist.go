package webpubsub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WebPubSubHub
}

type HubsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WebPubSubHub
}

type HubsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *HubsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// HubsList ...
func (c WebPubSubClient) HubsList(ctx context.Context, id WebPubSubId) (result HubsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &HubsListCustomPager{},
		Path:       fmt.Sprintf("%s/hubs", id.ID()),
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
		Values *[]WebPubSubHub `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// HubsListComplete retrieves all the results into a single object
func (c WebPubSubClient) HubsListComplete(ctx context.Context, id WebPubSubId) (HubsListCompleteResult, error) {
	return c.HubsListCompleteMatchingPredicate(ctx, id, WebPubSubHubOperationPredicate{})
}

// HubsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebPubSubClient) HubsListCompleteMatchingPredicate(ctx context.Context, id WebPubSubId, predicate WebPubSubHubOperationPredicate) (result HubsListCompleteResult, err error) {
	items := make([]WebPubSubHub, 0)

	resp, err := c.HubsList(ctx, id)
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

	result = HubsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
