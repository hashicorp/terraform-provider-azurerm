package productapilink

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProductApiLinkListByProductOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProductApiLinkContract
}

type WorkspaceProductApiLinkListByProductCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProductApiLinkContract
}

type WorkspaceProductApiLinkListByProductOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceProductApiLinkListByProductOperationOptions() WorkspaceProductApiLinkListByProductOperationOptions {
	return WorkspaceProductApiLinkListByProductOperationOptions{}
}

func (o WorkspaceProductApiLinkListByProductOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceProductApiLinkListByProductOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceProductApiLinkListByProductOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceProductApiLinkListByProductCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceProductApiLinkListByProductCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceProductApiLinkListByProduct ...
func (c ProductApiLinkClient) WorkspaceProductApiLinkListByProduct(ctx context.Context, id WorkspaceProductId, options WorkspaceProductApiLinkListByProductOperationOptions) (result WorkspaceProductApiLinkListByProductOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceProductApiLinkListByProductCustomPager{},
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
		Values *[]ProductApiLinkContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceProductApiLinkListByProductComplete retrieves all the results into a single object
func (c ProductApiLinkClient) WorkspaceProductApiLinkListByProductComplete(ctx context.Context, id WorkspaceProductId, options WorkspaceProductApiLinkListByProductOperationOptions) (WorkspaceProductApiLinkListByProductCompleteResult, error) {
	return c.WorkspaceProductApiLinkListByProductCompleteMatchingPredicate(ctx, id, options, ProductApiLinkContractOperationPredicate{})
}

// WorkspaceProductApiLinkListByProductCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProductApiLinkClient) WorkspaceProductApiLinkListByProductCompleteMatchingPredicate(ctx context.Context, id WorkspaceProductId, options WorkspaceProductApiLinkListByProductOperationOptions, predicate ProductApiLinkContractOperationPredicate) (result WorkspaceProductApiLinkListByProductCompleteResult, err error) {
	items := make([]ProductApiLinkContract, 0)

	resp, err := c.WorkspaceProductApiLinkListByProduct(ctx, id, options)
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

	result = WorkspaceProductApiLinkListByProductCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
