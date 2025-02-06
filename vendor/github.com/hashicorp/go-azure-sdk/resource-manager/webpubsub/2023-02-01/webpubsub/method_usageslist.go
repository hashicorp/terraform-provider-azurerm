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

type UsagesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SignalRServiceUsage
}

type UsagesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SignalRServiceUsage
}

type UsagesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *UsagesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// UsagesList ...
func (c WebPubSubClient) UsagesList(ctx context.Context, id LocationId) (result UsagesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &UsagesListCustomPager{},
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
		Values *[]SignalRServiceUsage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UsagesListComplete retrieves all the results into a single object
func (c WebPubSubClient) UsagesListComplete(ctx context.Context, id LocationId) (UsagesListCompleteResult, error) {
	return c.UsagesListCompleteMatchingPredicate(ctx, id, SignalRServiceUsageOperationPredicate{})
}

// UsagesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebPubSubClient) UsagesListCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate SignalRServiceUsageOperationPredicate) (result UsagesListCompleteResult, err error) {
	items := make([]SignalRServiceUsage, 0)

	resp, err := c.UsagesList(ctx, id)
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

	result = UsagesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
