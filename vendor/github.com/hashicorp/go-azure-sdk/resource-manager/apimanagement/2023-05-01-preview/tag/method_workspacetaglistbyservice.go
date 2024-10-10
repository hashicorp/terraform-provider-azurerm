package tag

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceTagListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagContract
}

type WorkspaceTagListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagContract
}

type WorkspaceTagListByServiceOperationOptions struct {
	Filter *string
	Scope  *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceTagListByServiceOperationOptions() WorkspaceTagListByServiceOperationOptions {
	return WorkspaceTagListByServiceOperationOptions{}
}

func (o WorkspaceTagListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceTagListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceTagListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Scope != nil {
		out.Append("scope", fmt.Sprintf("%v", *o.Scope))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceTagListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceTagListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceTagListByService ...
func (c TagClient) WorkspaceTagListByService(ctx context.Context, id WorkspaceId, options WorkspaceTagListByServiceOperationOptions) (result WorkspaceTagListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceTagListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/tags", id.ID()),
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
		Values *[]TagContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceTagListByServiceComplete retrieves all the results into a single object
func (c TagClient) WorkspaceTagListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspaceTagListByServiceOperationOptions) (WorkspaceTagListByServiceCompleteResult, error) {
	return c.WorkspaceTagListByServiceCompleteMatchingPredicate(ctx, id, options, TagContractOperationPredicate{})
}

// WorkspaceTagListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TagClient) WorkspaceTagListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceTagListByServiceOperationOptions, predicate TagContractOperationPredicate) (result WorkspaceTagListByServiceCompleteResult, err error) {
	items := make([]TagContract, 0)

	resp, err := c.WorkspaceTagListByService(ctx, id, options)
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

	result = WorkspaceTagListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
