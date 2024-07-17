package policyfragment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePolicyFragmentListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PolicyFragmentContract
}

type WorkspacePolicyFragmentListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PolicyFragmentContract
}

type WorkspacePolicyFragmentListByServiceOperationOptions struct {
	Filter  *string
	Orderby *string
	Skip    *int64
	Top     *int64
}

func DefaultWorkspacePolicyFragmentListByServiceOperationOptions() WorkspacePolicyFragmentListByServiceOperationOptions {
	return WorkspacePolicyFragmentListByServiceOperationOptions{}
}

func (o WorkspacePolicyFragmentListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspacePolicyFragmentListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspacePolicyFragmentListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspacePolicyFragmentListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspacePolicyFragmentListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspacePolicyFragmentListByService ...
func (c PolicyFragmentClient) WorkspacePolicyFragmentListByService(ctx context.Context, id WorkspaceId, options WorkspacePolicyFragmentListByServiceOperationOptions) (result WorkspacePolicyFragmentListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspacePolicyFragmentListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/policyFragments", id.ID()),
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
		Values *[]PolicyFragmentContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspacePolicyFragmentListByServiceComplete retrieves all the results into a single object
func (c PolicyFragmentClient) WorkspacePolicyFragmentListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspacePolicyFragmentListByServiceOperationOptions) (WorkspacePolicyFragmentListByServiceCompleteResult, error) {
	return c.WorkspacePolicyFragmentListByServiceCompleteMatchingPredicate(ctx, id, options, PolicyFragmentContractOperationPredicate{})
}

// WorkspacePolicyFragmentListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PolicyFragmentClient) WorkspacePolicyFragmentListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspacePolicyFragmentListByServiceOperationOptions, predicate PolicyFragmentContractOperationPredicate) (result WorkspacePolicyFragmentListByServiceCompleteResult, err error) {
	items := make([]PolicyFragmentContract, 0)

	resp, err := c.WorkspacePolicyFragmentListByService(ctx, id, options)
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

	result = WorkspacePolicyFragmentListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
