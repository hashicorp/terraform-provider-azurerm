package apischema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiSchemaListByApiOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SchemaContract
}

type WorkspaceApiSchemaListByApiCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SchemaContract
}

type WorkspaceApiSchemaListByApiOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceApiSchemaListByApiOperationOptions() WorkspaceApiSchemaListByApiOperationOptions {
	return WorkspaceApiSchemaListByApiOperationOptions{}
}

func (o WorkspaceApiSchemaListByApiOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiSchemaListByApiOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiSchemaListByApiOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceApiSchemaListByApiCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceApiSchemaListByApiCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceApiSchemaListByApi ...
func (c ApiSchemaClient) WorkspaceApiSchemaListByApi(ctx context.Context, id WorkspaceApiId, options WorkspaceApiSchemaListByApiOperationOptions) (result WorkspaceApiSchemaListByApiOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceApiSchemaListByApiCustomPager{},
		Path:          fmt.Sprintf("%s/schemas", id.ID()),
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
		Values *[]SchemaContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiSchemaListByApiComplete retrieves all the results into a single object
func (c ApiSchemaClient) WorkspaceApiSchemaListByApiComplete(ctx context.Context, id WorkspaceApiId, options WorkspaceApiSchemaListByApiOperationOptions) (WorkspaceApiSchemaListByApiCompleteResult, error) {
	return c.WorkspaceApiSchemaListByApiCompleteMatchingPredicate(ctx, id, options, SchemaContractOperationPredicate{})
}

// WorkspaceApiSchemaListByApiCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiSchemaClient) WorkspaceApiSchemaListByApiCompleteMatchingPredicate(ctx context.Context, id WorkspaceApiId, options WorkspaceApiSchemaListByApiOperationOptions, predicate SchemaContractOperationPredicate) (result WorkspaceApiSchemaListByApiCompleteResult, err error) {
	items := make([]SchemaContract, 0)

	resp, err := c.WorkspaceApiSchemaListByApi(ctx, id, options)
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

	result = WorkspaceApiSchemaListByApiCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
