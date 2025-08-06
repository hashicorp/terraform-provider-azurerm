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

type ListInstanceFunctionsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FunctionEnvelope
}

type ListInstanceFunctionsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FunctionEnvelope
}

type ListInstanceFunctionsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInstanceFunctionsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInstanceFunctionsSlot ...
func (c WebAppsClient) ListInstanceFunctionsSlot(ctx context.Context, id SlotId) (result ListInstanceFunctionsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInstanceFunctionsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/functions", id.ID()),
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
		Values *[]FunctionEnvelope `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInstanceFunctionsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListInstanceFunctionsSlotComplete(ctx context.Context, id SlotId) (ListInstanceFunctionsSlotCompleteResult, error) {
	return c.ListInstanceFunctionsSlotCompleteMatchingPredicate(ctx, id, FunctionEnvelopeOperationPredicate{})
}

// ListInstanceFunctionsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListInstanceFunctionsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate FunctionEnvelopeOperationPredicate) (result ListInstanceFunctionsSlotCompleteResult, err error) {
	items := make([]FunctionEnvelope, 0)

	resp, err := c.ListInstanceFunctionsSlot(ctx, id)
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

	result = ListInstanceFunctionsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
