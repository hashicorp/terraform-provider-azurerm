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

type WorkspacePolicyFragmentListReferencesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Resource
}

type WorkspacePolicyFragmentListReferencesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Resource
}

type WorkspacePolicyFragmentListReferencesOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultWorkspacePolicyFragmentListReferencesOperationOptions() WorkspacePolicyFragmentListReferencesOperationOptions {
	return WorkspacePolicyFragmentListReferencesOperationOptions{}
}

func (o WorkspacePolicyFragmentListReferencesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspacePolicyFragmentListReferencesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspacePolicyFragmentListReferencesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspacePolicyFragmentListReferencesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspacePolicyFragmentListReferencesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspacePolicyFragmentListReferences ...
func (c PolicyFragmentClient) WorkspacePolicyFragmentListReferences(ctx context.Context, id WorkspacePolicyFragmentId, options WorkspacePolicyFragmentListReferencesOperationOptions) (result WorkspacePolicyFragmentListReferencesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &WorkspacePolicyFragmentListReferencesCustomPager{},
		Path:          fmt.Sprintf("%s/listReferences", id.ID()),
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
		Values *[]Resource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspacePolicyFragmentListReferencesComplete retrieves all the results into a single object
func (c PolicyFragmentClient) WorkspacePolicyFragmentListReferencesComplete(ctx context.Context, id WorkspacePolicyFragmentId, options WorkspacePolicyFragmentListReferencesOperationOptions) (WorkspacePolicyFragmentListReferencesCompleteResult, error) {
	return c.WorkspacePolicyFragmentListReferencesCompleteMatchingPredicate(ctx, id, options, ResourceOperationPredicate{})
}

// WorkspacePolicyFragmentListReferencesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PolicyFragmentClient) WorkspacePolicyFragmentListReferencesCompleteMatchingPredicate(ctx context.Context, id WorkspacePolicyFragmentId, options WorkspacePolicyFragmentListReferencesOperationOptions, predicate ResourceOperationPredicate) (result WorkspacePolicyFragmentListReferencesCompleteResult, err error) {
	items := make([]Resource, 0)

	resp, err := c.WorkspacePolicyFragmentListReferences(ctx, id, options)
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

	result = WorkspacePolicyFragmentListReferencesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
