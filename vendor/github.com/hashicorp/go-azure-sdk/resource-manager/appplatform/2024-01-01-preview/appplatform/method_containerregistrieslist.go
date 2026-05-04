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

type ContainerRegistriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContainerRegistryResource
}

type ContainerRegistriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ContainerRegistryResource
}

type ContainerRegistriesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ContainerRegistriesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ContainerRegistriesList ...
func (c AppPlatformClient) ContainerRegistriesList(ctx context.Context, id commonids.SpringCloudServiceId) (result ContainerRegistriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ContainerRegistriesListCustomPager{},
		Path:       fmt.Sprintf("%s/containerRegistries", id.ID()),
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
		Values *[]ContainerRegistryResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ContainerRegistriesListComplete retrieves all the results into a single object
func (c AppPlatformClient) ContainerRegistriesListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ContainerRegistriesListCompleteResult, error) {
	return c.ContainerRegistriesListCompleteMatchingPredicate(ctx, id, ContainerRegistryResourceOperationPredicate{})
}

// ContainerRegistriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ContainerRegistriesListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate ContainerRegistryResourceOperationPredicate) (result ContainerRegistriesListCompleteResult, err error) {
	items := make([]ContainerRegistryResource, 0)

	resp, err := c.ContainerRegistriesList(ctx, id)
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

	result = ContainerRegistriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
