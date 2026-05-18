package sqlvirtualmachinegroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachinesListBySqlVMGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlVirtualMachine
}

type SqlVirtualMachinesListBySqlVMGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlVirtualMachine
}

type SqlVirtualMachinesListBySqlVMGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlVirtualMachinesListBySqlVMGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlVirtualMachinesListBySqlVMGroup ...
func (c SqlVirtualMachineGroupsClient) SqlVirtualMachinesListBySqlVMGroup(ctx context.Context, id SqlVirtualMachineGroupId) (result SqlVirtualMachinesListBySqlVMGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlVirtualMachinesListBySqlVMGroupCustomPager{},
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

// SqlVirtualMachinesListBySqlVMGroupComplete retrieves all the results into a single object
func (c SqlVirtualMachineGroupsClient) SqlVirtualMachinesListBySqlVMGroupComplete(ctx context.Context, id SqlVirtualMachineGroupId) (SqlVirtualMachinesListBySqlVMGroupCompleteResult, error) {
	return c.SqlVirtualMachinesListBySqlVMGroupCompleteMatchingPredicate(ctx, id, SqlVirtualMachineOperationPredicate{})
}

// SqlVirtualMachinesListBySqlVMGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SqlVirtualMachineGroupsClient) SqlVirtualMachinesListBySqlVMGroupCompleteMatchingPredicate(ctx context.Context, id SqlVirtualMachineGroupId, predicate SqlVirtualMachineOperationPredicate) (result SqlVirtualMachinesListBySqlVMGroupCompleteResult, err error) {
	items := make([]SqlVirtualMachine, 0)

	resp, err := c.SqlVirtualMachinesListBySqlVMGroup(ctx, id)
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

	result = SqlVirtualMachinesListBySqlVMGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
