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

type AvailablePrivateEndpointTypesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailablePrivateEndpointType
}

type AvailablePrivateEndpointTypesListCompleteResult struct {
	Items []AvailablePrivateEndpointType
}

// AvailablePrivateEndpointTypesList ...
func (c PrivateEndpointsClient) AvailablePrivateEndpointTypesList(ctx context.Context, id LocationId) (result AvailablePrivateEndpointTypesListOperationResponse, err error) {
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

// AvailablePrivateEndpointTypesListComplete retrieves all the results into a single object
func (c PrivateEndpointsClient) AvailablePrivateEndpointTypesListComplete(ctx context.Context, id LocationId) (AvailablePrivateEndpointTypesListCompleteResult, error) {
	return c.AvailablePrivateEndpointTypesListCompleteMatchingPredicate(ctx, id, AvailablePrivateEndpointTypeOperationPredicate{})
}

// AvailablePrivateEndpointTypesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateEndpointsClient) AvailablePrivateEndpointTypesListCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate AvailablePrivateEndpointTypeOperationPredicate) (result AvailablePrivateEndpointTypesListCompleteResult, err error) {
	items := make([]AvailablePrivateEndpointType, 0)

	resp, err := c.AvailablePrivateEndpointTypesList(ctx, id)
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

	result = AvailablePrivateEndpointTypesListCompleteResult{
		Items: items,
	}
	return
}
