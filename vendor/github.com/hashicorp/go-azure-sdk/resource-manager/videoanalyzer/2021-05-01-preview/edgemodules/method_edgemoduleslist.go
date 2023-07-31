package edgemodules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeModulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EdgeModuleEntity
}

type EdgeModulesListCompleteResult struct {
	Items []EdgeModuleEntity
}

type EdgeModulesListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultEdgeModulesListOperationOptions() EdgeModulesListOperationOptions {
	return EdgeModulesListOperationOptions{}
}

func (o EdgeModulesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o EdgeModulesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o EdgeModulesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// EdgeModulesList ...
func (c EdgeModulesClient) EdgeModulesList(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOperationOptions) (result EdgeModulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/edgeModules", id.ID()),
		OptionsObject: options,
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
		Values *[]EdgeModuleEntity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// EdgeModulesListComplete retrieves all the results into a single object
func (c EdgeModulesClient) EdgeModulesListComplete(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOperationOptions) (EdgeModulesListCompleteResult, error) {
	return c.EdgeModulesListCompleteMatchingPredicate(ctx, id, options, EdgeModuleEntityOperationPredicate{})
}

// EdgeModulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EdgeModulesClient) EdgeModulesListCompleteMatchingPredicate(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOperationOptions, predicate EdgeModuleEntityOperationPredicate) (result EdgeModulesListCompleteResult, err error) {
	items := make([]EdgeModuleEntity, 0)

	resp, err := c.EdgeModulesList(ctx, id, options)
	if err != nil {
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

	result = EdgeModulesListCompleteResult{
		Items: items,
	}
	return
}
