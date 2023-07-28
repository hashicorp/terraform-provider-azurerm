package modelversion

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
	Model        *[]ModelVersionResource
}

type ListCompleteResult struct {
	Items []ModelVersionResource
}

type ListOperationOptions struct {
	Description  *string
	Feed         *string
	ListViewType *ListViewType
	Offset       *int64
	OrderBy      *string
	Properties   *string
	Skip         *string
	Tags         *string
	Top          *int64
	Version      *string
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
	if o.Description != nil {
		out.Append("description", fmt.Sprintf("%v", *o.Description))
	}
	if o.Feed != nil {
		out.Append("feed", fmt.Sprintf("%v", *o.Feed))
	}
	if o.ListViewType != nil {
		out.Append("listViewType", fmt.Sprintf("%v", *o.ListViewType))
	}
	if o.Offset != nil {
		out.Append("offset", fmt.Sprintf("%v", *o.Offset))
	}
	if o.OrderBy != nil {
		out.Append("$orderBy", fmt.Sprintf("%v", *o.OrderBy))
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
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	if o.Version != nil {
		out.Append("version", fmt.Sprintf("%v", *o.Version))
	}
	return &out
}

// List ...
func (c ModelVersionClient) List(ctx context.Context, id ModelId, options ListOperationOptions) (result ListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/versions", id.ID()),
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
		Values *[]ModelVersionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListComplete retrieves all the results into a single object
func (c ModelVersionClient) ListComplete(ctx context.Context, id ModelId, options ListOperationOptions) (ListCompleteResult, error) {
	return c.ListCompleteMatchingPredicate(ctx, id, options, ModelVersionResourceOperationPredicate{})
}

// ListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ModelVersionClient) ListCompleteMatchingPredicate(ctx context.Context, id ModelId, options ListOperationOptions, predicate ModelVersionResourceOperationPredicate) (result ListCompleteResult, err error) {
	items := make([]ModelVersionResource, 0)

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
