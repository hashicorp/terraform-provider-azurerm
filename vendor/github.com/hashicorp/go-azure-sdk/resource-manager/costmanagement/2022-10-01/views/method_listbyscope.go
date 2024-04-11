package views

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

type ListByScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]View
}

type ListByScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []View
}

// ListByScope ...
func (c ViewsClient) ListByScope(ctx context.Context, id commonids.ScopeId) (result ListByScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.CostManagement/views", id.ID()),
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
		Values *[]View `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByScopeComplete retrieves all the results into a single object
func (c ViewsClient) ListByScopeComplete(ctx context.Context, id commonids.ScopeId) (ListByScopeCompleteResult, error) {
	return c.ListByScopeCompleteMatchingPredicate(ctx, id, ViewOperationPredicate{})
}

// ListByScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ViewsClient) ListByScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate ViewOperationPredicate) (result ListByScopeCompleteResult, err error) {
	items := make([]View, 0)

	resp, err := c.ListByScope(ctx, id)
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

	result = ListByScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
