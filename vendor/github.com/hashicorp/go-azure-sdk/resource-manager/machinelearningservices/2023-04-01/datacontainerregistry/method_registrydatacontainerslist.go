package datacontainerregistry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryDataContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataContainerResource
}

type RegistryDataContainersListCompleteResult struct {
	Items []DataContainerResource
}

type RegistryDataContainersListOperationOptions struct {
	ListViewType *ListViewType
	Skip         *string
}

func DefaultRegistryDataContainersListOperationOptions() RegistryDataContainersListOperationOptions {
	return RegistryDataContainersListOperationOptions{}
}

func (o RegistryDataContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryDataContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryDataContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ListViewType != nil {
		out.Append("listViewType", fmt.Sprintf("%v", *o.ListViewType))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// RegistryDataContainersList ...
func (c DataContainerRegistryClient) RegistryDataContainersList(ctx context.Context, id RegistryId, options RegistryDataContainersListOperationOptions) (result RegistryDataContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/data", id.ID()),
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
		Values *[]DataContainerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryDataContainersListComplete retrieves all the results into a single object
func (c DataContainerRegistryClient) RegistryDataContainersListComplete(ctx context.Context, id RegistryId, options RegistryDataContainersListOperationOptions) (RegistryDataContainersListCompleteResult, error) {
	return c.RegistryDataContainersListCompleteMatchingPredicate(ctx, id, options, DataContainerResourceOperationPredicate{})
}

// RegistryDataContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataContainerRegistryClient) RegistryDataContainersListCompleteMatchingPredicate(ctx context.Context, id RegistryId, options RegistryDataContainersListOperationOptions, predicate DataContainerResourceOperationPredicate) (result RegistryDataContainersListCompleteResult, err error) {
	items := make([]DataContainerResource, 0)

	resp, err := c.RegistryDataContainersList(ctx, id, options)
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

	result = RegistryDataContainersListCompleteResult{
		Items: items,
	}
	return
}
