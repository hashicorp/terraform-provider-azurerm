package registrymanagement

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

type RegistriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RegistryTrackedResource
}

type RegistriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RegistryTrackedResource
}

type RegistriesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RegistriesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RegistriesList ...
func (c RegistryManagementClient) RegistriesList(ctx context.Context, id commonids.ResourceGroupId) (result RegistriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RegistriesListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.MachineLearningServices/registries", id.ID()),
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
		Values *[]RegistryTrackedResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistriesListComplete retrieves all the results into a single object
func (c RegistryManagementClient) RegistriesListComplete(ctx context.Context, id commonids.ResourceGroupId) (RegistriesListCompleteResult, error) {
	return c.RegistriesListCompleteMatchingPredicate(ctx, id, RegistryTrackedResourceOperationPredicate{})
}

// RegistriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RegistryManagementClient) RegistriesListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate RegistryTrackedResourceOperationPredicate) (result RegistriesListCompleteResult, err error) {
	items := make([]RegistryTrackedResource, 0)

	resp, err := c.RegistriesList(ctx, id)
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

	result = RegistriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
