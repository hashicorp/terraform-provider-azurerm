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

type ListByClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Offer
}

type ListByClusterCompleteResult struct {
	Items []Offer
}

type ListByClusterOperationOptions struct {
	Expand *string
}

func DefaultListByClusterOperationOptions() ListByClusterOperationOptions {
	return ListByClusterOperationOptions{}
}

func (o ListByClusterOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByClusterOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByClusterOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// ListByCluster ...
func (c OffersClient) ListByCluster(ctx context.Context, id ClusterId, options ListByClusterOperationOptions) (result ListByClusterOperationResponse, err error) {
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

// ListByClusterComplete retrieves all the results into a single object
func (c OffersClient) ListByClusterComplete(ctx context.Context, id ClusterId, options ListByClusterOperationOptions) (ListByClusterCompleteResult, error) {
	return c.ListByClusterCompleteMatchingPredicate(ctx, id, options, OfferOperationPredicate{})
}

// ListByClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OffersClient) ListByClusterCompleteMatchingPredicate(ctx context.Context, id ClusterId, options ListByClusterOperationOptions, predicate OfferOperationPredicate) (result ListByClusterCompleteResult, err error) {
	items := make([]Offer, 0)

	resp, err := c.ListByCluster(ctx, id, options)
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

	result = ListByClusterCompleteResult{
		Items: items,
	}
	return
}
