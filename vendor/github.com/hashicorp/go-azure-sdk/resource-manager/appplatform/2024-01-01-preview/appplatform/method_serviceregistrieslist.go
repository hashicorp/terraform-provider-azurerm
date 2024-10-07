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

type ServiceRegistriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServiceRegistryResource
}

type ServiceRegistriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServiceRegistryResource
}

type ServiceRegistriesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServiceRegistriesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServiceRegistriesList ...
func (c AppPlatformClient) ServiceRegistriesList(ctx context.Context, id commonids.SpringCloudServiceId) (result ServiceRegistriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ServiceRegistriesListCustomPager{},
		Path:       fmt.Sprintf("%s/serviceRegistries", id.ID()),
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
		Values *[]ServiceRegistryResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServiceRegistriesListComplete retrieves all the results into a single object
func (c AppPlatformClient) ServiceRegistriesListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ServiceRegistriesListCompleteResult, error) {
	return c.ServiceRegistriesListCompleteMatchingPredicate(ctx, id, ServiceRegistryResourceOperationPredicate{})
}

// ServiceRegistriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ServiceRegistriesListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate ServiceRegistryResourceOperationPredicate) (result ServiceRegistriesListCompleteResult, err error) {
	items := make([]ServiceRegistryResource, 0)

	resp, err := c.ServiceRegistriesList(ctx, id)
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

	result = ServiceRegistriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
