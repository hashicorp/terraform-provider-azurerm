package tagproductlink

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceTagProductLinkListByProductOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagProductLinkContract
}

type WorkspaceTagProductLinkListByProductCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagProductLinkContract
}

type WorkspaceTagProductLinkListByProductOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceTagProductLinkListByProductOperationOptions() WorkspaceTagProductLinkListByProductOperationOptions {
	return WorkspaceTagProductLinkListByProductOperationOptions{}
}

func (o WorkspaceTagProductLinkListByProductOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceTagProductLinkListByProductOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceTagProductLinkListByProductOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceTagProductLinkListByProductCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceTagProductLinkListByProductCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceTagProductLinkListByProduct ...
func (c TagProductLinkClient) WorkspaceTagProductLinkListByProduct(ctx context.Context, id WorkspaceTagId, options WorkspaceTagProductLinkListByProductOperationOptions) (result WorkspaceTagProductLinkListByProductOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceTagProductLinkListByProductCustomPager{},
		Path:          fmt.Sprintf("%s/productLinks", id.ID()),
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
		Values *[]TagProductLinkContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceTagProductLinkListByProductComplete retrieves all the results into a single object
func (c TagProductLinkClient) WorkspaceTagProductLinkListByProductComplete(ctx context.Context, id WorkspaceTagId, options WorkspaceTagProductLinkListByProductOperationOptions) (WorkspaceTagProductLinkListByProductCompleteResult, error) {
	return c.WorkspaceTagProductLinkListByProductCompleteMatchingPredicate(ctx, id, options, TagProductLinkContractOperationPredicate{})
}

// WorkspaceTagProductLinkListByProductCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TagProductLinkClient) WorkspaceTagProductLinkListByProductCompleteMatchingPredicate(ctx context.Context, id WorkspaceTagId, options WorkspaceTagProductLinkListByProductOperationOptions, predicate TagProductLinkContractOperationPredicate) (result WorkspaceTagProductLinkListByProductCompleteResult, err error) {
	items := make([]TagProductLinkContract, 0)

	resp, err := c.WorkspaceTagProductLinkListByProduct(ctx, id, options)
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

	result = WorkspaceTagProductLinkListByProductCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
