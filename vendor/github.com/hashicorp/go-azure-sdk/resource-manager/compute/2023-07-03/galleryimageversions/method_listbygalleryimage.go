package galleryimageversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByGalleryImageOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GalleryImageVersion
}

type ListByGalleryImageCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GalleryImageVersion
}

type ListByGalleryImageCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByGalleryImageCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByGalleryImage ...
func (c GalleryImageVersionsClient) ListByGalleryImage(ctx context.Context, id GalleryImageId) (result ListByGalleryImageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByGalleryImageCustomPager{},
		Path:       fmt.Sprintf("%s/versions", id.ID()),
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
		Values *[]GalleryImageVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByGalleryImageComplete retrieves all the results into a single object
func (c GalleryImageVersionsClient) ListByGalleryImageComplete(ctx context.Context, id GalleryImageId) (ListByGalleryImageCompleteResult, error) {
	return c.ListByGalleryImageCompleteMatchingPredicate(ctx, id, GalleryImageVersionOperationPredicate{})
}

// ListByGalleryImageCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GalleryImageVersionsClient) ListByGalleryImageCompleteMatchingPredicate(ctx context.Context, id GalleryImageId, predicate GalleryImageVersionOperationPredicate) (result ListByGalleryImageCompleteResult, err error) {
	items := make([]GalleryImageVersion, 0)

	resp, err := c.ListByGalleryImage(ctx, id)
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

	result = ListByGalleryImageCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
