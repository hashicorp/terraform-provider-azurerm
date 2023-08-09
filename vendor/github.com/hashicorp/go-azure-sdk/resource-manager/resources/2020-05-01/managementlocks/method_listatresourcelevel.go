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

type ListAtResourceLevelOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagementLockObject
}

type ListAtResourceLevelCompleteResult struct {
	Items []ManagementLockObject
}

type ListAtResourceLevelOperationOptions struct {
	Filter *string
}

func DefaultListAtResourceLevelOperationOptions() ListAtResourceLevelOperationOptions {
	return ListAtResourceLevelOperationOptions{}
}

func (o ListAtResourceLevelOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAtResourceLevelOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListAtResourceLevelOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ListAtResourceLevel ...
func (c ManagementLocksClient) ListAtResourceLevel(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions) (result ListAtResourceLevelOperationResponse, err error) {
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

// ListAtResourceLevelComplete retrieves all the results into a single object
func (c ManagementLocksClient) ListAtResourceLevelComplete(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions) (ListAtResourceLevelCompleteResult, error) {
	return c.ListAtResourceLevelCompleteMatchingPredicate(ctx, id, options, ManagementLockObjectOperationPredicate{})
}

// ListAtResourceLevelCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagementLocksClient) ListAtResourceLevelCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions, predicate ManagementLockObjectOperationPredicate) (result ListAtResourceLevelCompleteResult, err error) {
	items := make([]ManagementLockObject, 0)

	resp, err := c.ListAtResourceLevel(ctx, id, options)
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

	result = ListAtResourceLevelCompleteResult{
		Items: items,
	}
	return
}
