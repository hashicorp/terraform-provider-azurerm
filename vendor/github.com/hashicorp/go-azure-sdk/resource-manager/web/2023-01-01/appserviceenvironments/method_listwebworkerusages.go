package appserviceenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListWebWorkerUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type ListWebWorkerUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type ListWebWorkerUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListWebWorkerUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListWebWorkerUsages ...
func (c AppServiceEnvironmentsClient) ListWebWorkerUsages(ctx context.Context, id WorkerPoolId) (result ListWebWorkerUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListWebWorkerUsagesCustomPager{},
		Path:       fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]Usage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListWebWorkerUsagesComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListWebWorkerUsagesComplete(ctx context.Context, id WorkerPoolId) (ListWebWorkerUsagesCompleteResult, error) {
	return c.ListWebWorkerUsagesCompleteMatchingPredicate(ctx, id, UsageOperationPredicate{})
}

// ListWebWorkerUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListWebWorkerUsagesCompleteMatchingPredicate(ctx context.Context, id WorkerPoolId, predicate UsageOperationPredicate) (result ListWebWorkerUsagesCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.ListWebWorkerUsages(ctx, id)
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

	result = ListWebWorkerUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
