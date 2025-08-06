package skus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByOfferOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Sku
}

type ListByOfferCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Sku
}

type ListByOfferOperationOptions struct {
	Expand *string
}

func DefaultListByOfferOperationOptions() ListByOfferOperationOptions {
	return ListByOfferOperationOptions{}
}

func (o ListByOfferOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByOfferOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByOfferOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

type ListByOfferCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByOfferCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByOffer ...
func (c SkusClient) ListByOffer(ctx context.Context, id OfferId, options ListByOfferOperationOptions) (result ListByOfferOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByOfferCustomPager{},
		Path:          fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]Sku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByOfferComplete retrieves all the results into a single object
func (c SkusClient) ListByOfferComplete(ctx context.Context, id OfferId, options ListByOfferOperationOptions) (ListByOfferCompleteResult, error) {
	return c.ListByOfferCompleteMatchingPredicate(ctx, id, options, SkuOperationPredicate{})
}

// ListByOfferCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SkusClient) ListByOfferCompleteMatchingPredicate(ctx context.Context, id OfferId, options ListByOfferOperationOptions, predicate SkuOperationPredicate) (result ListByOfferCompleteResult, err error) {
	items := make([]Sku, 0)

	resp, err := c.ListByOffer(ctx, id, options)
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

	result = ListByOfferCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
