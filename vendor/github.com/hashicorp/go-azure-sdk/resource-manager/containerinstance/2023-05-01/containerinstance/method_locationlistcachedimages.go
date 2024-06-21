package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationListCachedImagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CachedImages
}

type LocationListCachedImagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CachedImages
}

// LocationListCachedImages ...
func (c ContainerInstanceClient) LocationListCachedImages(ctx context.Context, id LocationId) (result LocationListCachedImagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/cachedImages", id.ID()),
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
		Values *[]CachedImages `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocationListCachedImagesComplete retrieves all the results into a single object
func (c ContainerInstanceClient) LocationListCachedImagesComplete(ctx context.Context, id LocationId) (LocationListCachedImagesCompleteResult, error) {
	return c.LocationListCachedImagesCompleteMatchingPredicate(ctx, id, CachedImagesOperationPredicate{})
}

// LocationListCachedImagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerInstanceClient) LocationListCachedImagesCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate CachedImagesOperationPredicate) (result LocationListCachedImagesCompleteResult, err error) {
	items := make([]CachedImages, 0)

	resp, err := c.LocationListCachedImages(ctx, id)
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

	result = LocationListCachedImagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
