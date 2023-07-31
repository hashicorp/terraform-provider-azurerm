package encodings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Transform
}

type TransformsListCompleteResult struct {
	Items []Transform
}

type TransformsListOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultTransformsListOperationOptions() TransformsListOperationOptions {
	return TransformsListOperationOptions{}
}

func (o TransformsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TransformsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o TransformsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	return &out
}

// TransformsList ...
func (c EncodingsClient) TransformsList(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions) (result TransformsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/transforms", id.ID()),
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
		Values *[]Transform `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TransformsListComplete retrieves all the results into a single object
func (c EncodingsClient) TransformsListComplete(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions) (TransformsListCompleteResult, error) {
	return c.TransformsListCompleteMatchingPredicate(ctx, id, options, TransformOperationPredicate{})
}

// TransformsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EncodingsClient) TransformsListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions, predicate TransformOperationPredicate) (result TransformsListCompleteResult, err error) {
	items := make([]Transform, 0)

	resp, err := c.TransformsList(ctx, id, options)
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

	result = TransformsListCompleteResult{
		Items: items,
	}
	return
}
