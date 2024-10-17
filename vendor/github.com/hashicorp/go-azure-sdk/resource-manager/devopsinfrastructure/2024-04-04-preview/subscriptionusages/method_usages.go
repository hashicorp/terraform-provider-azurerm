package subscriptionusages

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Quota
}

type UsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Quota
}

type UsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *UsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// Usages ...
func (c SubscriptionUsagesClient) Usages(ctx context.Context, id LocationId) (result UsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &UsagesCustomPager{},
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
		Values *[]Quota `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UsagesComplete retrieves all the results into a single object
func (c SubscriptionUsagesClient) UsagesComplete(ctx context.Context, id LocationId) (UsagesCompleteResult, error) {
	return c.UsagesCompleteMatchingPredicate(ctx, id, QuotaOperationPredicate{})
}

// UsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SubscriptionUsagesClient) UsagesCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate QuotaOperationPredicate) (result UsagesCompleteResult, err error) {
	items := make([]Quota, 0)

	resp, err := c.Usages(ctx, id)
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

	result = UsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
