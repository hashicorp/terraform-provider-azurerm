package tagapilink

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceTagApiLinkListByProductOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagApiLinkContract
}

type WorkspaceTagApiLinkListByProductCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagApiLinkContract
}

type WorkspaceTagApiLinkListByProductOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceTagApiLinkListByProductOperationOptions() WorkspaceTagApiLinkListByProductOperationOptions {
	return WorkspaceTagApiLinkListByProductOperationOptions{}
}

func (o WorkspaceTagApiLinkListByProductOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceTagApiLinkListByProductOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceTagApiLinkListByProductOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceTagApiLinkListByProductCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceTagApiLinkListByProductCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceTagApiLinkListByProduct ...
func (c TagApiLinkClient) WorkspaceTagApiLinkListByProduct(ctx context.Context, id WorkspaceTagId, options WorkspaceTagApiLinkListByProductOperationOptions) (result WorkspaceTagApiLinkListByProductOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceTagApiLinkListByProductCustomPager{},
		Path:          fmt.Sprintf("%s/apiLinks", id.ID()),
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
		Values *[]TagApiLinkContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceTagApiLinkListByProductComplete retrieves all the results into a single object
func (c TagApiLinkClient) WorkspaceTagApiLinkListByProductComplete(ctx context.Context, id WorkspaceTagId, options WorkspaceTagApiLinkListByProductOperationOptions) (WorkspaceTagApiLinkListByProductCompleteResult, error) {
	return c.WorkspaceTagApiLinkListByProductCompleteMatchingPredicate(ctx, id, options, TagApiLinkContractOperationPredicate{})
}

// WorkspaceTagApiLinkListByProductCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TagApiLinkClient) WorkspaceTagApiLinkListByProductCompleteMatchingPredicate(ctx context.Context, id WorkspaceTagId, options WorkspaceTagApiLinkListByProductOperationOptions, predicate TagApiLinkContractOperationPredicate) (result WorkspaceTagApiLinkListByProductCompleteResult, err error) {
	items := make([]TagApiLinkContract, 0)

	resp, err := c.WorkspaceTagApiLinkListByProduct(ctx, id, options)
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

	result = WorkspaceTagApiLinkListByProductCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
