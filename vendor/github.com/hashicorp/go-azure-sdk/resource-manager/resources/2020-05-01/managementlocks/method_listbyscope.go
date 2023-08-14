package managementlocks

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

type ListByScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagementLockObject
}

type ListByScopeCompleteResult struct {
	Items []ManagementLockObject
}

type ListByScopeOperationOptions struct {
	Filter *string
}

func DefaultListByScopeOperationOptions() ListByScopeOperationOptions {
	return ListByScopeOperationOptions{}
}

func (o ListByScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ListByScope ...
func (c ManagementLocksClient) ListByScope(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions) (result ListByScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/locks", id.ID()),
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
		Values *[]ManagementLockObject `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByScopeComplete retrieves all the results into a single object
func (c ManagementLocksClient) ListByScopeComplete(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions) (ListByScopeCompleteResult, error) {
	return c.ListByScopeCompleteMatchingPredicate(ctx, id, options, ManagementLockObjectOperationPredicate{})
}

// ListByScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagementLocksClient) ListByScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions, predicate ManagementLockObjectOperationPredicate) (result ListByScopeCompleteResult, err error) {
	items := make([]ManagementLockObject, 0)

	resp, err := c.ListByScope(ctx, id, options)
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

	result = ListByScopeCompleteResult{
		Items: items,
	}
	return
}
