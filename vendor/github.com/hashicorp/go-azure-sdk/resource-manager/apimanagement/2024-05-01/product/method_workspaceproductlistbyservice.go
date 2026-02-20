package product

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProductListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProductContract
}

type WorkspaceProductListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProductContract
}

type WorkspaceProductListByServiceOperationOptions struct {
	ExpandGroups *bool
	Filter       *string
	Skip         *int64
	Tags         *string
	Top          *int64
}

func DefaultWorkspaceProductListByServiceOperationOptions() WorkspaceProductListByServiceOperationOptions {
	return WorkspaceProductListByServiceOperationOptions{}
}

func (o WorkspaceProductListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceProductListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceProductListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ExpandGroups != nil {
		out.Append("expandGroups", fmt.Sprintf("%v", *o.ExpandGroups))
	}
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

type WorkspaceProductListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceProductListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceProductListByService ...
func (c ProductClient) WorkspaceProductListByService(ctx context.Context, id WorkspaceId, options WorkspaceProductListByServiceOperationOptions) (result WorkspaceProductListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceProductListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/products", id.ID()),
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
		Values *[]ProductContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceProductListByServiceComplete retrieves all the results into a single object
func (c ProductClient) WorkspaceProductListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspaceProductListByServiceOperationOptions) (WorkspaceProductListByServiceCompleteResult, error) {
	return c.WorkspaceProductListByServiceCompleteMatchingPredicate(ctx, id, options, ProductContractOperationPredicate{})
}

// WorkspaceProductListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProductClient) WorkspaceProductListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceProductListByServiceOperationOptions, predicate ProductContractOperationPredicate) (result WorkspaceProductListByServiceCompleteResult, err error) {
	items := make([]ProductContract, 0)

	resp, err := c.WorkspaceProductListByService(ctx, id, options)
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

	result = WorkspaceProductListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
