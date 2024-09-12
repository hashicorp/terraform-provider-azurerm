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

type ApiPortalsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiPortalResource
}

type ApiPortalsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiPortalResource
}

type ApiPortalsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApiPortalsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApiPortalsList ...
func (c AppPlatformClient) ApiPortalsList(ctx context.Context, id commonids.SpringCloudServiceId) (result ApiPortalsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ApiPortalsListCustomPager{},
		Path:       fmt.Sprintf("%s/apiPortals", id.ID()),
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
		Values *[]ApiPortalResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApiPortalsListComplete retrieves all the results into a single object
func (c AppPlatformClient) ApiPortalsListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ApiPortalsListCompleteResult, error) {
	return c.ApiPortalsListCompleteMatchingPredicate(ctx, id, ApiPortalResourceOperationPredicate{})
}

// ApiPortalsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ApiPortalsListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate ApiPortalResourceOperationPredicate) (result ApiPortalsListCompleteResult, err error) {
	items := make([]ApiPortalResource, 0)

	resp, err := c.ApiPortalsList(ctx, id)
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

	result = ApiPortalsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
