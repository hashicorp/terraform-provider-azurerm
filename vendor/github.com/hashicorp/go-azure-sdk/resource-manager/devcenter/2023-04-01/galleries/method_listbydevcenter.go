package galleries

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
	Model        *[]Gallery
}

type ListByDevCenterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Gallery
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
func (c GalleriesClient) ListByDevCenter(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions) (result ListByDevCenterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/galleries", id.ID()),
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
		Values *[]Gallery `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDevCenterComplete retrieves all the results into a single object
func (c GalleriesClient) ListByDevCenterComplete(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions) (ListByDevCenterCompleteResult, error) {
	return c.ListByDevCenterCompleteMatchingPredicate(ctx, id, options, GalleryOperationPredicate{})
}

// ListByDevCenterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GalleriesClient) ListByDevCenterCompleteMatchingPredicate(ctx context.Context, id DevCenterId, options ListByDevCenterOperationOptions, predicate GalleryOperationPredicate) (result ListByDevCenterCompleteResult, err error) {
	items := make([]Gallery, 0)

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
