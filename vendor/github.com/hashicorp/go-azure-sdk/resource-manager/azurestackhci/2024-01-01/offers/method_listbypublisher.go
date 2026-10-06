package offers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByPublisherOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Offer
}

type ListByPublisherCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Offer
}

type ListByPublisherOperationOptions struct {
	Expand *string
}

func DefaultListByPublisherOperationOptions() ListByPublisherOperationOptions {
	return ListByPublisherOperationOptions{}
}

func (o ListByPublisherOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByPublisherOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByPublisherOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

type ListByPublisherCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByPublisherCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByPublisher ...
func (c OffersClient) ListByPublisher(ctx context.Context, id PublisherId, options ListByPublisherOperationOptions) (result ListByPublisherOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByPublisherCustomPager{},
		Path:          fmt.Sprintf("%s/offers", id.ID()),
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
		Values *[]Offer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByPublisherComplete retrieves all the results into a single object
func (c OffersClient) ListByPublisherComplete(ctx context.Context, id PublisherId, options ListByPublisherOperationOptions) (ListByPublisherCompleteResult, error) {
	return c.ListByPublisherCompleteMatchingPredicate(ctx, id, options, OfferOperationPredicate{})
}

// ListByPublisherCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OffersClient) ListByPublisherCompleteMatchingPredicate(ctx context.Context, id PublisherId, options ListByPublisherOperationOptions, predicate OfferOperationPredicate) (result ListByPublisherCompleteResult, err error) {
	items := make([]Offer, 0)

	resp, err := c.ListByPublisher(ctx, id, options)
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

	result = ListByPublisherCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
