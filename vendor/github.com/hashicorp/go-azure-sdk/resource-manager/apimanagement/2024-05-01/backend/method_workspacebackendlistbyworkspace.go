package backend

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceBackendListByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BackendContract
}

type WorkspaceBackendListByWorkspaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BackendContract
}

type WorkspaceBackendListByWorkspaceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceBackendListByWorkspaceOperationOptions() WorkspaceBackendListByWorkspaceOperationOptions {
	return WorkspaceBackendListByWorkspaceOperationOptions{}
}

func (o WorkspaceBackendListByWorkspaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceBackendListByWorkspaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceBackendListByWorkspaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceBackendListByWorkspaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceBackendListByWorkspaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceBackendListByWorkspace ...
func (c BackendClient) WorkspaceBackendListByWorkspace(ctx context.Context, id WorkspaceId, options WorkspaceBackendListByWorkspaceOperationOptions) (result WorkspaceBackendListByWorkspaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceBackendListByWorkspaceCustomPager{},
		Path:          fmt.Sprintf("%s/backends", id.ID()),
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
		Values *[]BackendContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceBackendListByWorkspaceComplete retrieves all the results into a single object
func (c BackendClient) WorkspaceBackendListByWorkspaceComplete(ctx context.Context, id WorkspaceId, options WorkspaceBackendListByWorkspaceOperationOptions) (WorkspaceBackendListByWorkspaceCompleteResult, error) {
	return c.WorkspaceBackendListByWorkspaceCompleteMatchingPredicate(ctx, id, options, BackendContractOperationPredicate{})
}

// WorkspaceBackendListByWorkspaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BackendClient) WorkspaceBackendListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceBackendListByWorkspaceOperationOptions, predicate BackendContractOperationPredicate) (result WorkspaceBackendListByWorkspaceCompleteResult, err error) {
	items := make([]BackendContract, 0)

	resp, err := c.WorkspaceBackendListByWorkspace(ctx, id, options)
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

	result = WorkspaceBackendListByWorkspaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
