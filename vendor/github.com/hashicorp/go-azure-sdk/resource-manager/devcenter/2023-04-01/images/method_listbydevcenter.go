package images

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDevCenterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Image
}

type ListByDevCenterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Image
}

type ListByDevCenterOperationOptions struct {
	Top *int64
}

func DefaultListByDevCenterOperationOptions() ListByDevCenterOperationOptions {
	return ListByDevCenterOperationOptions{}
}

func (o ListByDevCenterOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByDevCenterOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByDevCenterOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByDevCenter ...
func (c ImagesClient) ListByDevCenter(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions) (result ListByDevCenterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/images", id.ID()),
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
		Values *[]Image `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDevCenterComplete retrieves all the results into a single object
func (c ImagesClient) ListByDevCenterComplete(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions) (ListByDevCenterCompleteResult, error) {
	return c.ListByDevCenterCompleteMatchingPredicate(ctx, id, options, ImageOperationPredicate{})
}

// ListByDevCenterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ImagesClient) ListByDevCenterCompleteMatchingPredicate(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions, predicate ImageOperationPredicate) (result ListByDevCenterCompleteResult, err error) {
	items := make([]Image, 0)

	resp, err := c.ListByDevCenter(ctx, id, options)
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

	result = ListByDevCenterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
