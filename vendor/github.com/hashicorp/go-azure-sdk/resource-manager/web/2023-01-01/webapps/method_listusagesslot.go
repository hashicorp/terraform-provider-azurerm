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

type ListUsagesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CsmUsageQuota
}

type ListUsagesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CsmUsageQuota
}

type ListUsagesSlotOperationOptions struct {
	Filter *string
}

func DefaultListUsagesSlotOperationOptions() ListUsagesSlotOperationOptions {
	return ListUsagesSlotOperationOptions{}
}

func (o ListUsagesSlotOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListUsagesSlotOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListUsagesSlotOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ListUsagesSlot ...
func (c WebAppsClient) ListUsagesSlot(ctx context.Context, id SlotId, options ListUsagesSlotOperationOptions) (result ListUsagesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/usages", id.ID()),
		OptionsObject: options,
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
		Values *[]CsmUsageQuota `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListUsagesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListUsagesSlotComplete(ctx context.Context, id SlotId, options ListUsagesSlotOperationOptions) (ListUsagesSlotCompleteResult, error) {
	return c.ListUsagesSlotCompleteMatchingPredicate(ctx, id, options, CsmUsageQuotaOperationPredicate{})
}

// ListUsagesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListUsagesSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, options ListUsagesSlotOperationOptions, predicate CsmUsageQuotaOperationPredicate) (result ListUsagesSlotCompleteResult, err error) {
	items := make([]CsmUsageQuota, 0)

	resp, err := c.ListUsagesSlot(ctx, id, options)
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

	result = ListUsagesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
