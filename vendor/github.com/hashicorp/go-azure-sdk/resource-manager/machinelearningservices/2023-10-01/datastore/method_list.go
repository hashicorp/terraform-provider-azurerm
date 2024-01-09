package datastore

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
	Model        *[]DatastoreResource
}

type ListCompleteResult struct {
	Items []DatastoreResource
}

type ListOperationOptions struct {
	Count      *int64
	IsDefault  *bool
	Names      *[]string
	OrderBy    *string
	OrderByAsc *bool
	SearchText *string
	Skip       *string
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
	if o.Count != nil {
		out.Append("count", fmt.Sprintf("%v", *o.Count))
	}
	if o.IsDefault != nil {
		out.Append("isDefault", fmt.Sprintf("%v", *o.IsDefault))
	}
	if o.Names != nil {
		out.Append("names", fmt.Sprintf("%v", *o.Names))
	}
	if o.OrderBy != nil {
		out.Append("orderBy", fmt.Sprintf("%v", *o.OrderBy))
	}
	if o.OrderByAsc != nil {
		out.Append("orderByAsc", fmt.Sprintf("%v", *o.OrderByAsc))
	}
	if o.SearchText != nil {
		out.Append("searchText", fmt.Sprintf("%v", *o.SearchText))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// List ...
func (c DatastoreClient) List(ctx context.Context, id WorkspaceId, options ListOperationOptions) (result ListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/dataStores", id.ID()),
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
		Values *[]DatastoreResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListComplete retrieves all the results into a single object
func (c DatastoreClient) ListComplete(ctx context.Context, id WorkspaceId, options ListOperationOptions) (ListCompleteResult, error) {
	return c.ListCompleteMatchingPredicate(ctx, id, options, DatastoreResourceOperationPredicate{})
}

// ListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DatastoreClient) ListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options ListOperationOptions, predicate DatastoreResourceOperationPredicate) (result ListCompleteResult, err error) {
	items := make([]DatastoreResource, 0)

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
