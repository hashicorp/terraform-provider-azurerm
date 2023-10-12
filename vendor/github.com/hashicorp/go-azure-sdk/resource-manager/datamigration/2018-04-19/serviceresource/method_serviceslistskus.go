package serviceresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesListSkusOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailableServiceSku
}

type ServicesListSkusCompleteResult struct {
	Items []AvailableServiceSku
}

// ServicesListSkus ...
func (c ServiceResourceClient) ServicesListSkus(ctx context.Context, id ServiceId) (result ServicesListSkusOperationResponse, err error) {
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
		Values *[]AvailableServiceSku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServicesListSkusComplete retrieves all the results into a single object
func (c ServiceResourceClient) ServicesListSkusComplete(ctx context.Context, id ServiceId) (ServicesListSkusCompleteResult, error) {
	return c.ServicesListSkusCompleteMatchingPredicate(ctx, id, AvailableServiceSkuOperationPredicate{})
}

// ServicesListSkusCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServiceResourceClient) ServicesListSkusCompleteMatchingPredicate(ctx context.Context, id ServiceId, predicate AvailableServiceSkuOperationPredicate) (result ServicesListSkusCompleteResult, err error) {
	items := make([]AvailableServiceSku, 0)

	resp, err := c.ServicesListSkus(ctx, id)
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

	result = ServicesListSkusCompleteResult{
		Items: items,
	}
	return
}
