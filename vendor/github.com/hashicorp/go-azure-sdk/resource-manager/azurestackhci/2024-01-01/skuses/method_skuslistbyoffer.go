package skuses

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusListByOfferOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Sku
}

type SkusListByOfferCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Sku
}

type SkusListByOfferOperationOptions struct {
	Expand *string
}

func DefaultSkusListByOfferOperationOptions() SkusListByOfferOperationOptions {
	return SkusListByOfferOperationOptions{}
}

func (o SkusListByOfferOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o SkusListByOfferOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SkusListByOfferOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

type SkusListByOfferCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SkusListByOfferCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SkusListByOffer ...
func (c SkusesClient) SkusListByOffer(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions) (result SkusListByOfferOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &SkusListByOfferCustomPager{},
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

// SkusListByOfferComplete retrieves all the results into a single object
func (c SkusesClient) SkusListByOfferComplete(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions) (SkusListByOfferCompleteResult, error) {
	return c.SkusListByOfferCompleteMatchingPredicate(ctx, id, options, SkuOperationPredicate{})
}

// SkusListByOfferCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SkusesClient) SkusListByOfferCompleteMatchingPredicate(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions, predicate SkuOperationPredicate) (result SkusListByOfferCompleteResult, err error) {
	items := make([]Sku, 0)

	resp, err := c.SkusListByOffer(ctx, id, options)
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

	result = SkusListByOfferCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
