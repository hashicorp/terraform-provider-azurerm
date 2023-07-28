package modelcontainer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryModelContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ModelContainerResource
}

type RegistryModelContainersListCompleteResult struct {
	Items []ModelContainerResource
}

type RegistryModelContainersListOperationOptions struct {
	ListViewType *ListViewType
	Skip         *string
}

func DefaultRegistryModelContainersListOperationOptions() RegistryModelContainersListOperationOptions {
	return RegistryModelContainersListOperationOptions{}
}

func (o RegistryModelContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryModelContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryModelContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ListViewType != nil {
		out.Append("listViewType", fmt.Sprintf("%v", *o.ListViewType))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// RegistryModelContainersList ...
func (c ModelContainerClient) RegistryModelContainersList(ctx context.Context, id RegistryId, options RegistryModelContainersListOperationOptions) (result RegistryModelContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/models", id.ID()),
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
		Values *[]ModelContainerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryModelContainersListComplete retrieves all the results into a single object
func (c ModelContainerClient) RegistryModelContainersListComplete(ctx context.Context, id RegistryId, options RegistryModelContainersListOperationOptions) (RegistryModelContainersListCompleteResult, error) {
	return c.RegistryModelContainersListCompleteMatchingPredicate(ctx, id, options, ModelContainerResourceOperationPredicate{})
}

// RegistryModelContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ModelContainerClient) RegistryModelContainersListCompleteMatchingPredicate(ctx context.Context, id RegistryId, options RegistryModelContainersListOperationOptions, predicate ModelContainerResourceOperationPredicate) (result RegistryModelContainersListCompleteResult, err error) {
	items := make([]ModelContainerResource, 0)

	resp, err := c.RegistryModelContainersList(ctx, id, options)
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

	result = RegistryModelContainersListCompleteResult{
		Items: items,
	}
	return
}
