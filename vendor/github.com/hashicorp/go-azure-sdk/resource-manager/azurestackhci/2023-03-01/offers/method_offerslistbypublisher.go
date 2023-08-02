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

type OffersListByPublisherOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Offer
}

type OffersListByPublisherCompleteResult struct {
	Items []Offer
}

type OffersListByPublisherOperationOptions struct {
	Expand *string
}

func DefaultOffersListByPublisherOperationOptions() OffersListByPublisherOperationOptions {
	return OffersListByPublisherOperationOptions{}
}

func (o OffersListByPublisherOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o OffersListByPublisherOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o OffersListByPublisherOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// OffersListByPublisher ...
func (c OffersClient) OffersListByPublisher(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions) (result OffersListByPublisherOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/offers", id.ID()),
		OptionsObject: options,
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

// OffersListByPublisherComplete retrieves all the results into a single object
func (c OffersClient) OffersListByPublisherComplete(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions) (OffersListByPublisherCompleteResult, error) {
	return c.OffersListByPublisherCompleteMatchingPredicate(ctx, id, options, OfferOperationPredicate{})
}

// OffersListByPublisherCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OffersClient) OffersListByPublisherCompleteMatchingPredicate(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions, predicate OfferOperationPredicate) (result OffersListByPublisherCompleteResult, err error) {
	items := make([]Offer, 0)

	resp, err := c.OffersListByPublisher(ctx, id, options)
	if err != nil {
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

	result = OffersListByPublisherCompleteResult{
		Items: items,
	}
	return
}
