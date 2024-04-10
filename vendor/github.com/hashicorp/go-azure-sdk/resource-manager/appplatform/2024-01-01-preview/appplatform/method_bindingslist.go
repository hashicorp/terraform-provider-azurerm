package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BindingsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BindingResource
}

type BindingsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BindingResource
}

// BindingsList ...
func (c AppPlatformClient) BindingsList(ctx context.Context, id AppId) (result BindingsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/bindings", id.ID()),
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
		Values *[]BindingResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BindingsListComplete retrieves all the results into a single object
func (c AppPlatformClient) BindingsListComplete(ctx context.Context, id AppId) (BindingsListCompleteResult, error) {
	return c.BindingsListCompleteMatchingPredicate(ctx, id, BindingResourceOperationPredicate{})
}

// BindingsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BindingsListCompleteMatchingPredicate(ctx context.Context, id AppId, predicate BindingResourceOperationPredicate) (result BindingsListCompleteResult, err error) {
	items := make([]BindingResource, 0)

	resp, err := c.BindingsList(ctx, id)
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

	result = BindingsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
