package environmentversion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryEnvironmentVersionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EnvironmentVersionResource
}

type RegistryEnvironmentVersionsListCompleteResult struct {
	Items []EnvironmentVersionResource
}

type RegistryEnvironmentVersionsListOperationOptions struct {
	ListViewType *ListViewType
	OrderBy      *string
	Skip         *string
	Top          *int64
}

func DefaultRegistryEnvironmentVersionsListOperationOptions() RegistryEnvironmentVersionsListOperationOptions {
	return RegistryEnvironmentVersionsListOperationOptions{}
}

func (o RegistryEnvironmentVersionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryEnvironmentVersionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryEnvironmentVersionsListOperationOptions) ToQuery() *client.QueryParams {
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
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RegistryEnvironmentVersionsList ...
func (c EnvironmentVersionClient) RegistryEnvironmentVersionsList(ctx context.Context, id RegistryEnvironmentId, options RegistryEnvironmentVersionsListOperationOptions) (result RegistryEnvironmentVersionsListOperationResponse, err error) {
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
		Values *[]EnvironmentVersionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryEnvironmentVersionsListComplete retrieves all the results into a single object
func (c EnvironmentVersionClient) RegistryEnvironmentVersionsListComplete(ctx context.Context, id RegistryEnvironmentId, options RegistryEnvironmentVersionsListOperationOptions) (RegistryEnvironmentVersionsListCompleteResult, error) {
	return c.RegistryEnvironmentVersionsListCompleteMatchingPredicate(ctx, id, options, EnvironmentVersionResourceOperationPredicate{})
}

// RegistryEnvironmentVersionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentVersionClient) RegistryEnvironmentVersionsListCompleteMatchingPredicate(ctx context.Context, id RegistryEnvironmentId, options RegistryEnvironmentVersionsListOperationOptions, predicate EnvironmentVersionResourceOperationPredicate) (result RegistryEnvironmentVersionsListCompleteResult, err error) {
	items := make([]EnvironmentVersionResource, 0)

	resp, err := c.RegistryEnvironmentVersionsList(ctx, id, options)
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

	result = RegistryEnvironmentVersionsListCompleteResult{
		Items: items,
	}
	return
}
