package sqlvirtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySqlVMGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlVirtualMachine
}

type ListBySqlVMGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlVirtualMachine
}

type ListBySqlVMGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBySqlVMGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBySqlVMGroup ...
func (c SqlVirtualMachinesClient) ListBySqlVMGroup(ctx context.Context, id SqlVirtualMachineGroupId) (result ListBySqlVMGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBySqlVMGroupCustomPager{},
		Path:       fmt.Sprintf("%s/sqlVirtualMachines", id.ID()),
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
		Values *[]SqlVirtualMachine `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySqlVMGroupComplete retrieves all the results into a single object
func (c SqlVirtualMachinesClient) ListBySqlVMGroupComplete(ctx context.Context, id SqlVirtualMachineGroupId) (ListBySqlVMGroupCompleteResult, error) {
	return c.ListBySqlVMGroupCompleteMatchingPredicate(ctx, id, SqlVirtualMachineOperationPredicate{})
}

// ListBySqlVMGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SqlVirtualMachinesClient) ListBySqlVMGroupCompleteMatchingPredicate(ctx context.Context, id SqlVirtualMachineGroupId, predicate SqlVirtualMachineOperationPredicate) (result ListBySqlVMGroupCompleteResult, err error) {
	items := make([]SqlVirtualMachine, 0)

	resp, err := c.ListBySqlVMGroup(ctx, id)
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

	result = ListBySqlVMGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
