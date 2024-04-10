package policyassignments

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

type ListForManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PolicyAssignment
}

type ListForManagementGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PolicyAssignment
}

type ListForManagementGroupOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListForManagementGroupOperationOptions() ListForManagementGroupOperationOptions {
	return ListForManagementGroupOperationOptions{}
}

func (o ListForManagementGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListForManagementGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListForManagementGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListForManagementGroup ...
func (c PolicyAssignmentsClient) ListForManagementGroup(ctx context.Context, id commonids.ManagementGroupId, options ListForManagementGroupOperationOptions) (result ListForManagementGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/policyAssignments", id.ID()),
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
		Values *[]PolicyAssignment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListForManagementGroupComplete retrieves all the results into a single object
func (c PolicyAssignmentsClient) ListForManagementGroupComplete(ctx context.Context, id commonids.ManagementGroupId, options ListForManagementGroupOperationOptions) (ListForManagementGroupCompleteResult, error) {
	return c.ListForManagementGroupCompleteMatchingPredicate(ctx, id, options, PolicyAssignmentOperationPredicate{})
}

// ListForManagementGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PolicyAssignmentsClient) ListForManagementGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ManagementGroupId, options ListForManagementGroupOperationOptions, predicate PolicyAssignmentOperationPredicate) (result ListForManagementGroupCompleteResult, err error) {
	items := make([]PolicyAssignment, 0)

	resp, err := c.ListForManagementGroup(ctx, id, options)
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

	result = ListForManagementGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
