package module

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PowerShell72ModuleListByAutomationAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Module
}

type PowerShell72ModuleListByAutomationAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Module
}

type PowerShell72ModuleListByAutomationAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PowerShell72ModuleListByAutomationAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PowerShell72ModuleListByAutomationAccount ...
func (c ModuleClient) PowerShell72ModuleListByAutomationAccount(ctx context.Context, id AutomationAccountId) (result PowerShell72ModuleListByAutomationAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PowerShell72ModuleListByAutomationAccountCustomPager{},
		Path:       fmt.Sprintf("%s/powerShell72Modules", id.ID()),
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

// PowerShell72ModuleListByAutomationAccountComplete retrieves all the results into a single object
func (c ModuleClient) PowerShell72ModuleListByAutomationAccountComplete(ctx context.Context, id AutomationAccountId) (PowerShell72ModuleListByAutomationAccountCompleteResult, error) {
	return c.PowerShell72ModuleListByAutomationAccountCompleteMatchingPredicate(ctx, id, ModuleOperationPredicate{})
}

// PowerShell72ModuleListByAutomationAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ModuleClient) PowerShell72ModuleListByAutomationAccountCompleteMatchingPredicate(ctx context.Context, id AutomationAccountId, predicate ModuleOperationPredicate) (result PowerShell72ModuleListByAutomationAccountCompleteResult, err error) {
	items := make([]Module, 0)

	resp, err := c.PowerShell72ModuleListByAutomationAccount(ctx, id)
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

	result = PowerShell72ModuleListByAutomationAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
