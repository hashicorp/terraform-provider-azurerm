package privateendpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailablePrivateEndpointTypesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailablePrivateEndpointType
}

type AvailablePrivateEndpointTypesListByResourceGroupCompleteResult struct {
	Items []AvailablePrivateEndpointType
}

// AvailablePrivateEndpointTypesListByResourceGroup ...
func (c PrivateEndpointsClient) AvailablePrivateEndpointTypesListByResourceGroup(ctx context.Context, id ProviderLocationId) (result AvailablePrivateEndpointTypesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/availablePrivateEndpointTypes", id.ID()),
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
		Values *[]AvailablePrivateEndpointType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AvailablePrivateEndpointTypesListByResourceGroupComplete retrieves all the results into a single object
func (c PrivateEndpointsClient) AvailablePrivateEndpointTypesListByResourceGroupComplete(ctx context.Context, id ProviderLocationId) (AvailablePrivateEndpointTypesListByResourceGroupCompleteResult, error) {
	return c.AvailablePrivateEndpointTypesListByResourceGroupCompleteMatchingPredicate(ctx, id, AvailablePrivateEndpointTypeOperationPredicate{})
}

// AvailablePrivateEndpointTypesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateEndpointsClient) AvailablePrivateEndpointTypesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ProviderLocationId, predicate AvailablePrivateEndpointTypeOperationPredicate) (result AvailablePrivateEndpointTypesListByResourceGroupCompleteResult, err error) {
	items := make([]AvailablePrivateEndpointType, 0)

	resp, err := c.AvailablePrivateEndpointTypesListByResourceGroup(ctx, id)
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

	result = AvailablePrivateEndpointTypesListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
