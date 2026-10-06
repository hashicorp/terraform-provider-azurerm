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

type DevToolPortalsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DevToolPortalResource
}

type DevToolPortalsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DevToolPortalResource
}

type DevToolPortalsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DevToolPortalsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DevToolPortalsList ...
func (c AppPlatformClient) DevToolPortalsList(ctx context.Context, id commonids.SpringCloudServiceId) (result DevToolPortalsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DevToolPortalsListCustomPager{},
		Path:       fmt.Sprintf("%s/devToolPortals", id.ID()),
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
		Values *[]DevToolPortalResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DevToolPortalsListComplete retrieves all the results into a single object
func (c AppPlatformClient) DevToolPortalsListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (DevToolPortalsListCompleteResult, error) {
	return c.DevToolPortalsListCompleteMatchingPredicate(ctx, id, DevToolPortalResourceOperationPredicate{})
}

// DevToolPortalsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) DevToolPortalsListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate DevToolPortalResourceOperationPredicate) (result DevToolPortalsListCompleteResult, err error) {
	items := make([]DevToolPortalResource, 0)

	resp, err := c.DevToolPortalsList(ctx, id)
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

	result = DevToolPortalsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
