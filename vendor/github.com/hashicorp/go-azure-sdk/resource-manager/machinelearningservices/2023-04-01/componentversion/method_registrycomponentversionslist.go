package componentversion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryComponentVersionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ComponentVersionResource
}

type RegistryComponentVersionsListCompleteResult struct {
	Items []ComponentVersionResource
}

type RegistryComponentVersionsListOperationOptions struct {
	OrderBy *string
	Skip    *string
	Top     *int64
}

func DefaultRegistryComponentVersionsListOperationOptions() RegistryComponentVersionsListOperationOptions {
	return RegistryComponentVersionsListOperationOptions{}
}

func (o RegistryComponentVersionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryComponentVersionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryComponentVersionsListOperationOptions) ToQuery() *client.QueryParams {
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

// RegistryComponentVersionsList ...
func (c ComponentVersionClient) RegistryComponentVersionsList(ctx context.Context, id RegistryComponentId, options RegistryComponentVersionsListOperationOptions) (result RegistryComponentVersionsListOperationResponse, err error) {
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
		Values *[]ComponentVersionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryComponentVersionsListComplete retrieves all the results into a single object
func (c ComponentVersionClient) RegistryComponentVersionsListComplete(ctx context.Context, id RegistryComponentId, options RegistryComponentVersionsListOperationOptions) (RegistryComponentVersionsListCompleteResult, error) {
	return c.RegistryComponentVersionsListCompleteMatchingPredicate(ctx, id, options, ComponentVersionResourceOperationPredicate{})
}

// RegistryComponentVersionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ComponentVersionClient) RegistryComponentVersionsListCompleteMatchingPredicate(ctx context.Context, id RegistryComponentId, options RegistryComponentVersionsListOperationOptions, predicate ComponentVersionResourceOperationPredicate) (result RegistryComponentVersionsListCompleteResult, err error) {
	items := make([]ComponentVersionResource, 0)

	resp, err := c.RegistryComponentVersionsList(ctx, id, options)
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

	result = RegistryComponentVersionsListCompleteResult{
		Items: items,
	}
	return
}
