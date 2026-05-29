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

type BuildServiceListBuildServicesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BuildService
}

type BuildServiceListBuildServicesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BuildService
}

type BuildServiceListBuildServicesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceListBuildServicesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceListBuildServices ...
func (c AppPlatformClient) BuildServiceListBuildServices(ctx context.Context, id commonids.SpringCloudServiceId) (result BuildServiceListBuildServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceListBuildServicesCustomPager{},
		Path:       fmt.Sprintf("%s/buildServices", id.ID()),
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
		Values *[]BuildService `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceListBuildServicesComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceListBuildServicesComplete(ctx context.Context, id commonids.SpringCloudServiceId) (BuildServiceListBuildServicesCompleteResult, error) {
	return c.BuildServiceListBuildServicesCompleteMatchingPredicate(ctx, id, BuildServiceOperationPredicate{})
}

// BuildServiceListBuildServicesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceListBuildServicesCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate BuildServiceOperationPredicate) (result BuildServiceListBuildServicesCompleteResult, err error) {
	items := make([]BuildService, 0)

	resp, err := c.BuildServiceListBuildServices(ctx, id)
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

	result = BuildServiceListBuildServicesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
