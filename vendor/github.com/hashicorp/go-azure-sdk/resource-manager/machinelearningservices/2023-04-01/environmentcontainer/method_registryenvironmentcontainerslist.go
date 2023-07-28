package environmentcontainer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryEnvironmentContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EnvironmentContainerResource
}

type RegistryEnvironmentContainersListCompleteResult struct {
	Items []EnvironmentContainerResource
}

type RegistryEnvironmentContainersListOperationOptions struct {
	ListViewType *ListViewType
	Skip         *string
}

func DefaultRegistryEnvironmentContainersListOperationOptions() RegistryEnvironmentContainersListOperationOptions {
	return RegistryEnvironmentContainersListOperationOptions{}
}

func (o RegistryEnvironmentContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryEnvironmentContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryEnvironmentContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ListViewType != nil {
		out.Append("listViewType", fmt.Sprintf("%v", *o.ListViewType))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// RegistryEnvironmentContainersList ...
func (c EnvironmentContainerClient) RegistryEnvironmentContainersList(ctx context.Context, id RegistryId, options RegistryEnvironmentContainersListOperationOptions) (result RegistryEnvironmentContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/environments", id.ID()),
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
		Values *[]EnvironmentContainerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryEnvironmentContainersListComplete retrieves all the results into a single object
func (c EnvironmentContainerClient) RegistryEnvironmentContainersListComplete(ctx context.Context, id RegistryId, options RegistryEnvironmentContainersListOperationOptions) (RegistryEnvironmentContainersListCompleteResult, error) {
	return c.RegistryEnvironmentContainersListCompleteMatchingPredicate(ctx, id, options, EnvironmentContainerResourceOperationPredicate{})
}

// RegistryEnvironmentContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentContainerClient) RegistryEnvironmentContainersListCompleteMatchingPredicate(ctx context.Context, id RegistryId, options RegistryEnvironmentContainersListOperationOptions, predicate EnvironmentContainerResourceOperationPredicate) (result RegistryEnvironmentContainersListCompleteResult, err error) {
	items := make([]EnvironmentContainerResource, 0)

	resp, err := c.RegistryEnvironmentContainersList(ctx, id, options)
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

	result = RegistryEnvironmentContainersListCompleteResult{
		Items: items,
	}
	return
}
