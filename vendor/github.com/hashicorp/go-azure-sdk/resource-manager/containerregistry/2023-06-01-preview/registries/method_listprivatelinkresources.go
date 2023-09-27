package registries

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPrivateLinkResourcesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateLinkResource
}

type ListPrivateLinkResourcesCompleteResult struct {
	Items []PrivateLinkResource
}

// ListPrivateLinkResources ...
func (c RegistriesClient) ListPrivateLinkResources(ctx context.Context, id RegistryId) (result ListPrivateLinkResourcesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/privateLinkResources", id.ID()),
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
		Values *[]PrivateLinkResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListPrivateLinkResourcesComplete retrieves all the results into a single object
func (c RegistriesClient) ListPrivateLinkResourcesComplete(ctx context.Context, id RegistryId) (ListPrivateLinkResourcesCompleteResult, error) {
	return c.ListPrivateLinkResourcesCompleteMatchingPredicate(ctx, id, PrivateLinkResourceOperationPredicate{})
}

// ListPrivateLinkResourcesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RegistriesClient) ListPrivateLinkResourcesCompleteMatchingPredicate(ctx context.Context, id RegistryId, predicate PrivateLinkResourceOperationPredicate) (result ListPrivateLinkResourcesCompleteResult, err error) {
	items := make([]PrivateLinkResource, 0)

	resp, err := c.ListPrivateLinkResources(ctx, id)
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

	result = ListPrivateLinkResourcesCompleteResult{
		Items: items,
	}
	return
}
