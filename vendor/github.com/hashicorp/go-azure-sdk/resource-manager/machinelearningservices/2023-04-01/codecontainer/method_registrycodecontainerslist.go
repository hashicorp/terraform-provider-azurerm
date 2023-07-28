package codecontainer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryCodeContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CodeContainerResource
}

type RegistryCodeContainersListCompleteResult struct {
	Items []CodeContainerResource
}

type RegistryCodeContainersListOperationOptions struct {
	Skip *string
}

func DefaultRegistryCodeContainersListOperationOptions() RegistryCodeContainersListOperationOptions {
	return RegistryCodeContainersListOperationOptions{}
}

func (o RegistryCodeContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegistryCodeContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegistryCodeContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// RegistryCodeContainersList ...
func (c CodeContainerClient) RegistryCodeContainersList(ctx context.Context, id RegistryId, options RegistryCodeContainersListOperationOptions) (result RegistryCodeContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/codes", id.ID()),
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
		Values *[]CodeContainerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RegistryCodeContainersListComplete retrieves all the results into a single object
func (c CodeContainerClient) RegistryCodeContainersListComplete(ctx context.Context, id RegistryId, options RegistryCodeContainersListOperationOptions) (RegistryCodeContainersListCompleteResult, error) {
	return c.RegistryCodeContainersListCompleteMatchingPredicate(ctx, id, options, CodeContainerResourceOperationPredicate{})
}

// RegistryCodeContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CodeContainerClient) RegistryCodeContainersListCompleteMatchingPredicate(ctx context.Context, id RegistryId, options RegistryCodeContainersListOperationOptions, predicate CodeContainerResourceOperationPredicate) (result RegistryCodeContainersListCompleteResult, err error) {
	items := make([]CodeContainerResource, 0)

	resp, err := c.RegistryCodeContainersList(ctx, id, options)
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

	result = RegistryCodeContainersListCompleteResult{
		Items: items,
	}
	return
}
