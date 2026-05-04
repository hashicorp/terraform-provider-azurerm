package guestagents

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

type ListByVirtualMachineInstanceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GuestAgent
}

type ListByVirtualMachineInstanceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GuestAgent
}

type ListByVirtualMachineInstanceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByVirtualMachineInstanceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByVirtualMachineInstance ...
func (c GuestAgentsClient) ListByVirtualMachineInstance(ctx context.Context, id commonids.ScopeId) (result ListByVirtualMachineInstanceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByVirtualMachineInstanceCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.ScVmm/virtualMachineInstances/default/guestAgents", id.ID()),
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
		Values *[]GuestAgent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByVirtualMachineInstanceComplete retrieves all the results into a single object
func (c GuestAgentsClient) ListByVirtualMachineInstanceComplete(ctx context.Context, id commonids.ScopeId) (ListByVirtualMachineInstanceCompleteResult, error) {
	return c.ListByVirtualMachineInstanceCompleteMatchingPredicate(ctx, id, GuestAgentOperationPredicate{})
}

// ListByVirtualMachineInstanceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GuestAgentsClient) ListByVirtualMachineInstanceCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate GuestAgentOperationPredicate) (result ListByVirtualMachineInstanceCompleteResult, err error) {
	items := make([]GuestAgent, 0)

	resp, err := c.ListByVirtualMachineInstance(ctx, id)
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

	result = ListByVirtualMachineInstanceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
