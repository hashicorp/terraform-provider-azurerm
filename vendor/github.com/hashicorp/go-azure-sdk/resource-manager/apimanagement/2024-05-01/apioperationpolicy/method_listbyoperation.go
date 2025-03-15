package apioperationpolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByOperationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PolicyContract
}

type ListByOperationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PolicyContract
}

type ListByOperationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByOperationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByOperation ...
func (c ApiOperationPolicyClient) ListByOperation(ctx context.Context, id OperationId) (result ListByOperationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByOperationCustomPager{},
		Path:       fmt.Sprintf("%s/policies", id.ID()),
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
		Values *[]PolicyContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByOperationComplete retrieves all the results into a single object
func (c ApiOperationPolicyClient) ListByOperationComplete(ctx context.Context, id OperationId) (ListByOperationCompleteResult, error) {
	return c.ListByOperationCompleteMatchingPredicate(ctx, id, PolicyContractOperationPredicate{})
}

// ListByOperationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiOperationPolicyClient) ListByOperationCompleteMatchingPredicate(ctx context.Context, id OperationId, predicate PolicyContractOperationPredicate) (result ListByOperationCompleteResult, err error) {
	items := make([]PolicyContract, 0)

	resp, err := c.ListByOperation(ctx, id)
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

	result = ListByOperationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
