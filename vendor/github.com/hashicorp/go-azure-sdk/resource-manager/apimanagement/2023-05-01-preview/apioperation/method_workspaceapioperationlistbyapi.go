package apioperation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiOperationListByApiOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OperationContract
}

type WorkspaceApiOperationListByApiCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OperationContract
}

type WorkspaceApiOperationListByApiOperationOptions struct {
	Filter *string
	Skip   *int64
	Tags   *string
	Top    *int64
}

func DefaultWorkspaceApiOperationListByApiOperationOptions() WorkspaceApiOperationListByApiOperationOptions {
	return WorkspaceApiOperationListByApiOperationOptions{}
}

func (o WorkspaceApiOperationListByApiOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiOperationListByApiOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceApiOperationListByApiOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Tags != nil {
		out.Append("tags", fmt.Sprintf("%v", *o.Tags))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// WorkspaceApiOperationListByApi ...
func (c ApiOperationClient) WorkspaceApiOperationListByApi(ctx context.Context, id WorkspaceApiId, options WorkspaceApiOperationListByApiOperationOptions) (result WorkspaceApiOperationListByApiOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/operations", id.ID()),
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
		Values *[]OperationContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiOperationListByApiComplete retrieves all the results into a single object
func (c ApiOperationClient) WorkspaceApiOperationListByApiComplete(ctx context.Context, id WorkspaceApiId, options WorkspaceApiOperationListByApiOperationOptions) (WorkspaceApiOperationListByApiCompleteResult, error) {
	return c.WorkspaceApiOperationListByApiCompleteMatchingPredicate(ctx, id, options, OperationContractOperationPredicate{})
}

// WorkspaceApiOperationListByApiCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiOperationClient) WorkspaceApiOperationListByApiCompleteMatchingPredicate(ctx context.Context, id WorkspaceApiId, options WorkspaceApiOperationListByApiOperationOptions, predicate OperationContractOperationPredicate) (result WorkspaceApiOperationListByApiCompleteResult, err error) {
	items := make([]OperationContract, 0)

	resp, err := c.WorkspaceApiOperationListByApi(ctx, id, options)
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

	result = WorkspaceApiOperationListByApiCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
