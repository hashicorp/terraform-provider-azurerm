package diagnostic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceDiagnosticListByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DiagnosticContract
}

type WorkspaceDiagnosticListByWorkspaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DiagnosticContract
}

type WorkspaceDiagnosticListByWorkspaceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceDiagnosticListByWorkspaceOperationOptions() WorkspaceDiagnosticListByWorkspaceOperationOptions {
	return WorkspaceDiagnosticListByWorkspaceOperationOptions{}
}

func (o WorkspaceDiagnosticListByWorkspaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceDiagnosticListByWorkspaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceDiagnosticListByWorkspaceOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceDiagnosticListByWorkspaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceDiagnosticListByWorkspaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceDiagnosticListByWorkspace ...
func (c DiagnosticClient) WorkspaceDiagnosticListByWorkspace(ctx context.Context, id WorkspaceId, options WorkspaceDiagnosticListByWorkspaceOperationOptions) (result WorkspaceDiagnosticListByWorkspaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceDiagnosticListByWorkspaceCustomPager{},
		Path:          fmt.Sprintf("%s/diagnostics", id.ID()),
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
		Values *[]DiagnosticContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceDiagnosticListByWorkspaceComplete retrieves all the results into a single object
func (c DiagnosticClient) WorkspaceDiagnosticListByWorkspaceComplete(ctx context.Context, id WorkspaceId, options WorkspaceDiagnosticListByWorkspaceOperationOptions) (WorkspaceDiagnosticListByWorkspaceCompleteResult, error) {
	return c.WorkspaceDiagnosticListByWorkspaceCompleteMatchingPredicate(ctx, id, options, DiagnosticContractOperationPredicate{})
}

// WorkspaceDiagnosticListByWorkspaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DiagnosticClient) WorkspaceDiagnosticListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceDiagnosticListByWorkspaceOperationOptions, predicate DiagnosticContractOperationPredicate) (result WorkspaceDiagnosticListByWorkspaceCompleteResult, err error) {
	items := make([]DiagnosticContract, 0)

	resp, err := c.WorkspaceDiagnosticListByWorkspace(ctx, id, options)
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

	result = WorkspaceDiagnosticListByWorkspaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
