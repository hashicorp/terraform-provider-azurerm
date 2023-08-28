package galleryapplicationversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByGalleryApplicationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GalleryApplicationVersion
}

type ListByGalleryApplicationCompleteResult struct {
	Items []GalleryApplicationVersion
}

// ListByGalleryApplication ...
func (c GalleryApplicationVersionsClient) ListByGalleryApplication(ctx context.Context, id ApplicationId) (result ListByGalleryApplicationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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
		Values *[]GalleryApplicationVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByGalleryApplicationComplete retrieves all the results into a single object
func (c GalleryApplicationVersionsClient) ListByGalleryApplicationComplete(ctx context.Context, id ApplicationId) (ListByGalleryApplicationCompleteResult, error) {
	return c.ListByGalleryApplicationCompleteMatchingPredicate(ctx, id, GalleryApplicationVersionOperationPredicate{})
}

// ListByGalleryApplicationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GalleryApplicationVersionsClient) ListByGalleryApplicationCompleteMatchingPredicate(ctx context.Context, id ApplicationId, predicate GalleryApplicationVersionOperationPredicate) (result ListByGalleryApplicationCompleteResult, err error) {
	items := make([]GalleryApplicationVersion, 0)

	resp, err := c.ListByGalleryApplication(ctx, id)
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

	result = ListByGalleryApplicationCompleteResult{
		Items: items,
	}
	return
}
