package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetPrivateEndpointConnectionListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemotePrivateEndpointConnectionARMResource
}

type GetPrivateEndpointConnectionListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RemotePrivateEndpointConnectionARMResource
}

// GetPrivateEndpointConnectionList ...
func (c StaticSitesClient) GetPrivateEndpointConnectionList(ctx context.Context, id StaticSiteId) (result GetPrivateEndpointConnectionListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/privateEndpointConnections", id.ID()),
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
		Values *[]RemotePrivateEndpointConnectionARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetPrivateEndpointConnectionListComplete retrieves all the results into a single object
func (c StaticSitesClient) GetPrivateEndpointConnectionListComplete(ctx context.Context, id StaticSiteId) (GetPrivateEndpointConnectionListCompleteResult, error) {
	return c.GetPrivateEndpointConnectionListCompleteMatchingPredicate(ctx, id, RemotePrivateEndpointConnectionARMResourceOperationPredicate{})
}

// GetPrivateEndpointConnectionListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetPrivateEndpointConnectionListCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate RemotePrivateEndpointConnectionARMResourceOperationPredicate) (result GetPrivateEndpointConnectionListCompleteResult, err error) {
	items := make([]RemotePrivateEndpointConnectionARMResource, 0)

	resp, err := c.GetPrivateEndpointConnectionList(ctx, id)
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

	result = GetPrivateEndpointConnectionListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
