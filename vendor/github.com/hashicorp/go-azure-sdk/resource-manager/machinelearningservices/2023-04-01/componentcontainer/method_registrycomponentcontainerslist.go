package componentcontainer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryComponentContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ComponentContainerResource
}

type RegistryComponentContainersListCompleteResult struct {
	Items []ComponentContainerResource
}

type RegistryComponentContainersListOperationOptions struct {
	Skip *string
}

func DefaultRegistryComponentContainersListOperationOptions() RegistryComponentContainersListOperationOptions {
	return RegistryComponentContainersListOperationOptions{}
}

func (o RegistryComponentContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryComponentContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryComponentContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// RegistryComponentContainersList ...
func (c ComponentContainerClient) RegistryComponentContainersList(ctx context.Context, id RegistryId, options RegistryComponentContainersListOperationOptions) (result RegistryComponentContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/components", id.ID()),
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
		Values *[]ComponentContainerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryComponentContainersListComplete retrieves all the results into a single object
func (c ComponentContainerClient) RegistryComponentContainersListComplete(ctx context.Context, id RegistryId, options RegistryComponentContainersListOperationOptions) (RegistryComponentContainersListCompleteResult, error) {
	return c.RegistryComponentContainersListCompleteMatchingPredicate(ctx, id, options, ComponentContainerResourceOperationPredicate{})
}

// RegistryComponentContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ComponentContainerClient) RegistryComponentContainersListCompleteMatchingPredicate(ctx context.Context, id RegistryId, options RegistryComponentContainersListOperationOptions, predicate ComponentContainerResourceOperationPredicate) (result RegistryComponentContainersListCompleteResult, err error) {
	items := make([]ComponentContainerResource, 0)

	resp, err := c.RegistryComponentContainersList(ctx, id, options)
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

	result = RegistryComponentContainersListCompleteResult{
		Items: items,
	}
	return
}
