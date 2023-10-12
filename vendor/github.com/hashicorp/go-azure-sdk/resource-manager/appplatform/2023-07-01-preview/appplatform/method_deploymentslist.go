package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentResource
}

type DeploymentsListCompleteResult struct {
	Items []DeploymentResource
}

type DeploymentsListOperationOptions struct {
	Version *[]string
}

func DefaultDeploymentsListOperationOptions() DeploymentsListOperationOptions {
	return DeploymentsListOperationOptions{}
}

func (o DeploymentsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeploymentsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DeploymentsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Version != nil {
		out.Append("version", fmt.Sprintf("%v", *o.Version))
	}
	return &out
}

// DeploymentsList ...
func (c AppPlatformClient) DeploymentsList(ctx context.Context, id AppId, options DeploymentsListOperationOptions) (result DeploymentsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/deployments", id.ID()),
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
		Values *[]DeploymentResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DeploymentsListComplete retrieves all the results into a single object
func (c AppPlatformClient) DeploymentsListComplete(ctx context.Context, id AppId, options DeploymentsListOperationOptions) (DeploymentsListCompleteResult, error) {
	return c.DeploymentsListCompleteMatchingPredicate(ctx, id, options, DeploymentResourceOperationPredicate{})
}

// DeploymentsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) DeploymentsListCompleteMatchingPredicate(ctx context.Context, id AppId, options DeploymentsListOperationOptions, predicate DeploymentResourceOperationPredicate) (result DeploymentsListCompleteResult, err error) {
	items := make([]DeploymentResource, 0)

	resp, err := c.DeploymentsList(ctx, id, options)
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

	result = DeploymentsListCompleteResult{
		Items: items,
	}
	return
}
