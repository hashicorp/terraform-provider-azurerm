package webapps

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

type ListSlotsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Site
}

type ListSlotsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Site
}

// ListSlots ...
func (c WebAppsClient) ListSlots(ctx context.Context, id commonids.AppServiceId) (result ListSlotsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/slots", id.ID()),
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
		Values *[]Site `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSlotsComplete retrieves all the results into a single object
func (c WebAppsClient) ListSlotsComplete(ctx context.Context, id commonids.AppServiceId) (ListSlotsCompleteResult, error) {
	return c.ListSlotsCompleteMatchingPredicate(ctx, id, SiteOperationPredicate{})
}

// ListSlotsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSlotsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate SiteOperationPredicate) (result ListSlotsCompleteResult, err error) {
	items := make([]Site, 0)

	resp, err := c.ListSlots(ctx, id)
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

	result = ListSlotsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
