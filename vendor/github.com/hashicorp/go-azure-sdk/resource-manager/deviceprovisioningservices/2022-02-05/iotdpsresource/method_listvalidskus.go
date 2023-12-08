package iotdpsresource

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

type ListValidSkusOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]IotDpsSkuDefinition
}

type ListValidSkusCompleteResult struct {
	Items []IotDpsSkuDefinition
}

// ListValidSkus ...
func (c IotDpsResourceClient) ListValidSkus(ctx context.Context, id commonids.ProvisioningServiceId) (result ListValidSkusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]IotDpsSkuDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListValidSkusComplete retrieves all the results into a single object
func (c IotDpsResourceClient) ListValidSkusComplete(ctx context.Context, id commonids.ProvisioningServiceId) (ListValidSkusCompleteResult, error) {
	return c.ListValidSkusCompleteMatchingPredicate(ctx, id, IotDpsSkuDefinitionOperationPredicate{})
}

// ListValidSkusCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c IotDpsResourceClient) ListValidSkusCompleteMatchingPredicate(ctx context.Context, id commonids.ProvisioningServiceId, predicate IotDpsSkuDefinitionOperationPredicate) (result ListValidSkusCompleteResult, err error) {
	items := make([]IotDpsSkuDefinition, 0)

	resp, err := c.ListValidSkus(ctx, id)
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

	result = ListValidSkusCompleteResult{
		Items: items,
	}
	return
}
