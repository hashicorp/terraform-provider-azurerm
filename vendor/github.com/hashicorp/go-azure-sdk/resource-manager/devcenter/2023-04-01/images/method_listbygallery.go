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

type ListByGalleryOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Image
}

type ListByGalleryCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Image
}

type ListByGalleryOperationOptions struct {
	Top *int64
}

func DefaultListByGalleryOperationOptions() ListByGalleryOperationOptions {
	return ListByGalleryOperationOptions{}
}

func (o ListByGalleryOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByGalleryOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByGalleryOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByGalleryCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByGalleryCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByGallery ...
func (c ImagesClient) ListByGallery(ctx context.Context, id GalleryId, options ListByGalleryOperationOptions) (result ListByGalleryOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByGalleryCustomPager{},
		Path:          fmt.Sprintf("%s/images", id.ID()),
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

// ListByGalleryComplete retrieves all the results into a single object
func (c ImagesClient) ListByGalleryComplete(ctx context.Context, id GalleryId, options ListByGalleryOperationOptions) (ListByGalleryCompleteResult, error) {
	return c.ListByGalleryCompleteMatchingPredicate(ctx, id, options, ImageOperationPredicate{})
}

// ListByGalleryCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ImagesClient) ListByGalleryCompleteMatchingPredicate(ctx context.Context, id GalleryId, options ListByGalleryOperationOptions, predicate ImageOperationPredicate) (result ListByGalleryCompleteResult, err error) {
	items := make([]Image, 0)

	resp, err := c.ListByGallery(ctx, id, options)
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

	result = ListByGalleryCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
