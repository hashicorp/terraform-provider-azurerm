package policysetdefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBuiltInOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PolicySetDefinition
}

type ListBuiltInCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PolicySetDefinition
}

type ListBuiltInOperationOptions struct {
	Expand *string
	Filter *string
	Top    *int64
}

func DefaultListBuiltInOperationOptions() ListBuiltInOperationOptions {
	return ListBuiltInOperationOptions{}
}

func (o ListBuiltInOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListBuiltInOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListBuiltInOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListBuiltInCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBuiltInCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBuiltIn ...
func (c PolicySetDefinitionsClient) ListBuiltIn(ctx context.Context, options ListBuiltInOperationOptions) (result ListBuiltInOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListBuiltInCustomPager{},
		Path:          "/providers/Microsoft.Authorization/policySetDefinitions",
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
		Values *[]PolicySetDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBuiltInComplete retrieves all the results into a single object
func (c PolicySetDefinitionsClient) ListBuiltInComplete(ctx context.Context, options ListBuiltInOperationOptions) (ListBuiltInCompleteResult, error) {
	return c.ListBuiltInCompleteMatchingPredicate(ctx, options, PolicySetDefinitionOperationPredicate{})
}

// ListBuiltInCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PolicySetDefinitionsClient) ListBuiltInCompleteMatchingPredicate(ctx context.Context, options ListBuiltInOperationOptions, predicate PolicySetDefinitionOperationPredicate) (result ListBuiltInCompleteResult, err error) {
	items := make([]PolicySetDefinition, 0)

	resp, err := c.ListBuiltIn(ctx, options)
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

	result = ListBuiltInCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
