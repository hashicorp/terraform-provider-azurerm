package iotconnectors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FhirDestinationsListByIotConnectorOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]IotFhirDestination
}

type FhirDestinationsListByIotConnectorCompleteResult struct {
	Items []IotFhirDestination
}

// FhirDestinationsListByIotConnector ...
func (c IotConnectorsClient) FhirDestinationsListByIotConnector(ctx context.Context, id IotConnectorId) (result FhirDestinationsListByIotConnectorOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/fhirDestinations", id.ID()),
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
		Values *[]IotFhirDestination `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FhirDestinationsListByIotConnectorComplete retrieves all the results into a single object
func (c IotConnectorsClient) FhirDestinationsListByIotConnectorComplete(ctx context.Context, id IotConnectorId) (FhirDestinationsListByIotConnectorCompleteResult, error) {
	return c.FhirDestinationsListByIotConnectorCompleteMatchingPredicate(ctx, id, IotFhirDestinationOperationPredicate{})
}

// FhirDestinationsListByIotConnectorCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c IotConnectorsClient) FhirDestinationsListByIotConnectorCompleteMatchingPredicate(ctx context.Context, id IotConnectorId, predicate IotFhirDestinationOperationPredicate) (result FhirDestinationsListByIotConnectorCompleteResult, err error) {
	items := make([]IotFhirDestination, 0)

	resp, err := c.FhirDestinationsListByIotConnector(ctx, id)
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

	result = FhirDestinationsListByIotConnectorCompleteResult{
		Items: items,
	}
	return
}
