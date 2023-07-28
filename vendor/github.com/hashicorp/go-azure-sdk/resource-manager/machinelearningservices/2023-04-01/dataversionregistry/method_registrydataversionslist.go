package dataversionregistry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryDataVersionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataVersionBaseResource
}

type RegistryDataVersionsListCompleteResult struct {
	Items []DataVersionBaseResource
}

type RegistryDataVersionsListOperationOptions struct {
	ListViewType *ListViewType
	OrderBy      *string
	Skip         *string
	Tags         *string
	Top          *int64
}

func DefaultRegistryDataVersionsListOperationOptions() RegistryDataVersionsListOperationOptions {
	return RegistryDataVersionsListOperationOptions{}
}

func (o RegistryDataVersionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryDataVersionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryDataVersionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ListViewType != nil {
		out.Append("listViewType", fmt.Sprintf("%v", *o.ListViewType))
	}
	if o.OrderBy != nil {
		out.Append("$orderBy", fmt.Sprintf("%v", *o.OrderBy))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Tags != nil {
		out.Append("$tags", fmt.Sprintf("%v", *o.Tags))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RegistryDataVersionsList ...
func (c DataVersionRegistryClient) RegistryDataVersionsList(ctx context.Context, id DataId, options RegistryDataVersionsListOperationOptions) (result RegistryDataVersionsListOperationResponse, err error) {
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
		Values *[]DataVersionBaseResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryDataVersionsListComplete retrieves all the results into a single object
func (c DataVersionRegistryClient) RegistryDataVersionsListComplete(ctx context.Context, id DataId, options RegistryDataVersionsListOperationOptions) (RegistryDataVersionsListCompleteResult, error) {
	return c.RegistryDataVersionsListCompleteMatchingPredicate(ctx, id, options, DataVersionBaseResourceOperationPredicate{})
}

// RegistryDataVersionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataVersionRegistryClient) RegistryDataVersionsListCompleteMatchingPredicate(ctx context.Context, id DataId, options RegistryDataVersionsListOperationOptions, predicate DataVersionBaseResourceOperationPredicate) (result RegistryDataVersionsListCompleteResult, err error) {
	items := make([]DataVersionBaseResource, 0)

	resp, err := c.RegistryDataVersionsList(ctx, id, options)
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

	result = RegistryDataVersionsListCompleteResult{
		Items: items,
	}
	return
}
