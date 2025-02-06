package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StoragesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageResource
}

type StoragesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageResource
}

type StoragesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *StoragesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// StoragesList ...
func (c AppPlatformClient) StoragesList(ctx context.Context, id commonids.SpringCloudServiceId) (result StoragesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &StoragesListCustomPager{},
		Path:       fmt.Sprintf("%s/storages", id.ID()),
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
		Values *[]StorageResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// StoragesListComplete retrieves all the results into a single object
func (c AppPlatformClient) StoragesListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (StoragesListCompleteResult, error) {
	return c.StoragesListCompleteMatchingPredicate(ctx, id, StorageResourceOperationPredicate{})
}

// StoragesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) StoragesListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate StorageResourceOperationPredicate) (result StoragesListCompleteResult, err error) {
	items := make([]StorageResource, 0)

	resp, err := c.StoragesList(ctx, id)
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

	result = StoragesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
