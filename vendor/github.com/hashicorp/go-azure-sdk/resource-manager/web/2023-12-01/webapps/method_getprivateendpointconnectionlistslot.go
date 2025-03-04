package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetPrivateEndpointConnectionListSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemotePrivateEndpointConnectionARMResource
}

type GetPrivateEndpointConnectionListSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RemotePrivateEndpointConnectionARMResource
}

type GetPrivateEndpointConnectionListSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetPrivateEndpointConnectionListSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetPrivateEndpointConnectionListSlot ...
func (c WebAppsClient) GetPrivateEndpointConnectionListSlot(ctx context.Context, id SlotId) (result GetPrivateEndpointConnectionListSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetPrivateEndpointConnectionListSlotCustomPager{},
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

// GetPrivateEndpointConnectionListSlotComplete retrieves all the results into a single object
func (c WebAppsClient) GetPrivateEndpointConnectionListSlotComplete(ctx context.Context, id SlotId) (GetPrivateEndpointConnectionListSlotCompleteResult, error) {
	return c.GetPrivateEndpointConnectionListSlotCompleteMatchingPredicate(ctx, id, RemotePrivateEndpointConnectionARMResourceOperationPredicate{})
}

// GetPrivateEndpointConnectionListSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) GetPrivateEndpointConnectionListSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate RemotePrivateEndpointConnectionARMResourceOperationPredicate) (result GetPrivateEndpointConnectionListSlotCompleteResult, err error) {
	items := make([]RemotePrivateEndpointConnectionARMResource, 0)

	resp, err := c.GetPrivateEndpointConnectionListSlot(ctx, id)
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

	result = GetPrivateEndpointConnectionListSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
