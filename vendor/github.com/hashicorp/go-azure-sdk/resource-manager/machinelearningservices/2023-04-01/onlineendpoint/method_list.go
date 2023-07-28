package onlineendpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OnlineEndpointTrackedResource
}

type ListCompleteResult struct {
	Items []OnlineEndpointTrackedResource
}

type ListOperationOptions struct {
	ComputeType *EndpointComputeType
	Count       *int64
	Name        *string
	OrderBy     *OrderString
	Properties  *string
	Skip        *string
	Tags        *string
}

func DefaultListOperationOptions() ListOperationOptions {
	return ListOperationOptions{}
}

func (o ListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ComputeType != nil {
		out.Append("computeType", fmt.Sprintf("%v", *o.ComputeType))
	}
	if o.Count != nil {
		out.Append("count", fmt.Sprintf("%v", *o.Count))
	}
	if o.Name != nil {
		out.Append("name", fmt.Sprintf("%v", *o.Name))
	}
	if o.OrderBy != nil {
		out.Append("orderBy", fmt.Sprintf("%v", *o.OrderBy))
	}
	if o.Properties != nil {
		out.Append("properties", fmt.Sprintf("%v", *o.Properties))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Tags != nil {
		out.Append("tags", fmt.Sprintf("%v", *o.Tags))
	}
	return &out
}

// List ...
func (c OnlineEndpointClient) List(ctx context.Context, id WorkspaceId, options ListOperationOptions) (result ListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/onlineEndpoints", id.ID()),
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
		Values *[]OnlineEndpointTrackedResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListComplete retrieves all the results into a single object
func (c OnlineEndpointClient) ListComplete(ctx context.Context, id WorkspaceId, options ListOperationOptions) (ListCompleteResult, error) {
	return c.ListCompleteMatchingPredicate(ctx, id, options, OnlineEndpointTrackedResourceOperationPredicate{})
}

// ListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OnlineEndpointClient) ListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options ListOperationOptions, predicate OnlineEndpointTrackedResourceOperationPredicate) (result ListCompleteResult, err error) {
	items := make([]OnlineEndpointTrackedResource, 0)

	resp, err := c.List(ctx, id, options)
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

	result = ListCompleteResult{
		Items: items,
	}
	return
}
