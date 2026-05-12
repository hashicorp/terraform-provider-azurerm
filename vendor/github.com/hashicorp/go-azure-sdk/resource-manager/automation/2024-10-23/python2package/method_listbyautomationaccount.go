package python2package

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAutomationAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Module
}

type ListByAutomationAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Module
}

type ListByAutomationAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByAutomationAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByAutomationAccount ...
func (c Python2PackageClient) ListByAutomationAccount(ctx context.Context, id AutomationAccountId) (result ListByAutomationAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByAutomationAccountCustomPager{},
		Path:       fmt.Sprintf("%s/python2Packages", id.ID()),
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
		Values *[]Module `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAutomationAccountComplete retrieves all the results into a single object
func (c Python2PackageClient) ListByAutomationAccountComplete(ctx context.Context, id AutomationAccountId) (ListByAutomationAccountCompleteResult, error) {
	return c.ListByAutomationAccountCompleteMatchingPredicate(ctx, id, ModuleOperationPredicate{})
}

// ListByAutomationAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c Python2PackageClient) ListByAutomationAccountCompleteMatchingPredicate(ctx context.Context, id AutomationAccountId, predicate ModuleOperationPredicate) (result ListByAutomationAccountCompleteResult, err error) {
	items := make([]Module, 0)

	resp, err := c.ListByAutomationAccount(ctx, id)
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

	result = ListByAutomationAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
