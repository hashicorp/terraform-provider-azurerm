package codeversion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryCodeVersionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CodeVersionResource
}

type RegistryCodeVersionsListCompleteResult struct {
	Items []CodeVersionResource
}

type RegistryCodeVersionsListOperationOptions struct {
	OrderBy *string
	Skip    *string
	Top     *int64
}

func DefaultRegistryCodeVersionsListOperationOptions() RegistryCodeVersionsListOperationOptions {
	return RegistryCodeVersionsListOperationOptions{}
}

func (o RegistryCodeVersionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryCodeVersionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryCodeVersionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.OrderBy != nil {
		out.Append("$orderBy", fmt.Sprintf("%v", *o.OrderBy))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RegistryCodeVersionsList ...
func (c CodeVersionClient) RegistryCodeVersionsList(ctx context.Context, id RegistryCodeId, options RegistryCodeVersionsListOperationOptions) (result RegistryCodeVersionsListOperationResponse, err error) {
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
		Values *[]CodeVersionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryCodeVersionsListComplete retrieves all the results into a single object
func (c CodeVersionClient) RegistryCodeVersionsListComplete(ctx context.Context, id RegistryCodeId, options RegistryCodeVersionsListOperationOptions) (RegistryCodeVersionsListCompleteResult, error) {
	return c.RegistryCodeVersionsListCompleteMatchingPredicate(ctx, id, options, CodeVersionResourceOperationPredicate{})
}

// RegistryCodeVersionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CodeVersionClient) RegistryCodeVersionsListCompleteMatchingPredicate(ctx context.Context, id RegistryCodeId, options RegistryCodeVersionsListOperationOptions, predicate CodeVersionResourceOperationPredicate) (result RegistryCodeVersionsListCompleteResult, err error) {
	items := make([]CodeVersionResource, 0)

	resp, err := c.RegistryCodeVersionsList(ctx, id, options)
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

	result = RegistryCodeVersionsListCompleteResult{
		Items: items,
	}
	return
}
