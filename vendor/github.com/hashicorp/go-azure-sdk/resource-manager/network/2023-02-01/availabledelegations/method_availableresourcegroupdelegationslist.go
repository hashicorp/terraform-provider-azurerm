package availabledelegations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableResourceGroupDelegationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailableDelegation
}

type AvailableResourceGroupDelegationsListCompleteResult struct {
	Items []AvailableDelegation
}

// AvailableResourceGroupDelegationsList ...
func (c AvailableDelegationsClient) AvailableResourceGroupDelegationsList(ctx context.Context, id ProviderLocationId) (result AvailableResourceGroupDelegationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/availableDelegations", id.ID()),
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
		Values *[]AvailableDelegation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AvailableResourceGroupDelegationsListComplete retrieves all the results into a single object
func (c AvailableDelegationsClient) AvailableResourceGroupDelegationsListComplete(ctx context.Context, id ProviderLocationId) (AvailableResourceGroupDelegationsListCompleteResult, error) {
	return c.AvailableResourceGroupDelegationsListCompleteMatchingPredicate(ctx, id, AvailableDelegationOperationPredicate{})
}

// AvailableResourceGroupDelegationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AvailableDelegationsClient) AvailableResourceGroupDelegationsListCompleteMatchingPredicate(ctx context.Context, id ProviderLocationId, predicate AvailableDelegationOperationPredicate) (result AvailableResourceGroupDelegationsListCompleteResult, err error) {
	items := make([]AvailableDelegation, 0)

	resp, err := c.AvailableResourceGroupDelegationsList(ctx, id)
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

	result = AvailableResourceGroupDelegationsListCompleteResult{
		Items: items,
	}
	return
}
