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

type RegistriesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RegistryTrackedResource
}

type RegistriesListBySubscriptionCompleteResult struct {
	Items []RegistryTrackedResource
}

// RegistriesListBySubscription ...
func (c RegistryManagementClient) RegistriesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result RegistriesListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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

// RegistriesListBySubscriptionComplete retrieves all the results into a single object
func (c RegistryManagementClient) RegistriesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (RegistriesListBySubscriptionCompleteResult, error) {
	return c.RegistriesListBySubscriptionCompleteMatchingPredicate(ctx, id, RegistryTrackedResourceOperationPredicate{})
}

// RegistriesListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RegistryManagementClient) RegistriesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate RegistryTrackedResourceOperationPredicate) (result RegistriesListBySubscriptionCompleteResult, err error) {
	items := make([]RegistryTrackedResource, 0)

	resp, err := c.RegistriesListBySubscription(ctx, id)
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

	result = RegistriesListBySubscriptionCompleteResult{
		Items: items,
	}
	return
}
