package imageversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByImageOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ImageVersion
}

type ListByImageCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ImageVersion
}

type ListByImageCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByImageCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByImage ...
func (c ImageVersionsClient) ListByImage(ctx context.Context, id GalleryImageId) (result ListByImageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByImageCustomPager{},
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
		Values *[]ImageVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByImageComplete retrieves all the results into a single object
func (c ImageVersionsClient) ListByImageComplete(ctx context.Context, id GalleryImageId) (ListByImageCompleteResult, error) {
	return c.ListByImageCompleteMatchingPredicate(ctx, id, ImageVersionOperationPredicate{})
}

// ListByImageCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ImageVersionsClient) ListByImageCompleteMatchingPredicate(ctx context.Context, id GalleryImageId, predicate ImageVersionOperationPredicate) (result ListByImageCompleteResult, err error) {
	items := make([]ImageVersion, 0)

	resp, err := c.ListByImage(ctx, id)
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

	result = ListByImageCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
