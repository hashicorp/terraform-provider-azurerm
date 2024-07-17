package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiContract
}

type WorkspaceApiListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiContract
}

type WorkspaceApiListByServiceOperationOptions struct {
	ExpandApiVersionSet *bool
	Filter              *string
	Skip                *int64
	Tags                *string
	Top                 *int64
}

func DefaultWorkspaceApiListByServiceOperationOptions() WorkspaceApiListByServiceOperationOptions {
	return WorkspaceApiListByServiceOperationOptions{}
}

func (o WorkspaceApiListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceApiListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ExpandApiVersionSet != nil {
		out.Append("expandApiVersionSet", fmt.Sprintf("%v", *o.ExpandApiVersionSet))
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

type WorkspaceApiListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceApiListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceApiListByService ...
func (c ApiClient) WorkspaceApiListByService(ctx context.Context, id WorkspaceId, options WorkspaceApiListByServiceOperationOptions) (result WorkspaceApiListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceApiListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/apis", id.ID()),
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
		Values *[]ApiContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiListByServiceComplete retrieves all the results into a single object
func (c ApiClient) WorkspaceApiListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspaceApiListByServiceOperationOptions) (WorkspaceApiListByServiceCompleteResult, error) {
	return c.WorkspaceApiListByServiceCompleteMatchingPredicate(ctx, id, options, ApiContractOperationPredicate{})
}

// WorkspaceApiListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiClient) WorkspaceApiListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceApiListByServiceOperationOptions, predicate ApiContractOperationPredicate) (result WorkspaceApiListByServiceCompleteResult, err error) {
	items := make([]ApiContract, 0)

	resp, err := c.WorkspaceApiListByService(ctx, id, options)
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

	result = WorkspaceApiListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
