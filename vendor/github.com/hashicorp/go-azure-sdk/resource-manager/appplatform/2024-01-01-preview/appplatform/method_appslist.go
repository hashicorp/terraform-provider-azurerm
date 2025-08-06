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

type AppsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AppResource
}

type AppsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AppResource
}

type AppsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AppsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AppsList ...
func (c AppPlatformClient) AppsList(ctx context.Context, id commonids.SpringCloudServiceId) (result AppsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AppsListCustomPager{},
		Path:       fmt.Sprintf("%s/apps", id.ID()),
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
		Values *[]AppResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AppsListComplete retrieves all the results into a single object
func (c AppPlatformClient) AppsListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (AppsListCompleteResult, error) {
	return c.AppsListCompleteMatchingPredicate(ctx, id, AppResourceOperationPredicate{})
}

// AppsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) AppsListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate AppResourceOperationPredicate) (result AppsListCompleteResult, err error) {
	items := make([]AppResource, 0)

	resp, err := c.AppsList(ctx, id)
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

	result = AppsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
