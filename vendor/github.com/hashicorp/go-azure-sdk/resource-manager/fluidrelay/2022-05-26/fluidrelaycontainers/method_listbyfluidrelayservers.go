package fluidrelaycontainers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByFluidRelayServersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FluidRelayContainer
}

type ListByFluidRelayServersCompleteResult struct {
	Items []FluidRelayContainer
}

// ListByFluidRelayServers ...
func (c FluidRelayContainersClient) ListByFluidRelayServers(ctx context.Context, id FluidRelayServerId) (result ListByFluidRelayServersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/fluidRelayContainers", id.ID()),
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
		Values *[]FluidRelayContainer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByFluidRelayServersComplete retrieves all the results into a single object
func (c FluidRelayContainersClient) ListByFluidRelayServersComplete(ctx context.Context, id FluidRelayServerId) (ListByFluidRelayServersCompleteResult, error) {
	return c.ListByFluidRelayServersCompleteMatchingPredicate(ctx, id, FluidRelayContainerOperationPredicate{})
}

// ListByFluidRelayServersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FluidRelayContainersClient) ListByFluidRelayServersCompleteMatchingPredicate(ctx context.Context, id FluidRelayServerId, predicate FluidRelayContainerOperationPredicate) (result ListByFluidRelayServersCompleteResult, err error) {
	items := make([]FluidRelayContainer, 0)

	resp, err := c.ListByFluidRelayServers(ctx, id)
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

	result = ListByFluidRelayServersCompleteResult{
		Items: items,
	}
	return
}
