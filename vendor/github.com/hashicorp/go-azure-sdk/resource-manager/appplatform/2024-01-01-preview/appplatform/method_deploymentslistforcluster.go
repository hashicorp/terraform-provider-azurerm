package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentsListForClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentResource
}

type DeploymentsListForClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentResource
}

type DeploymentsListForClusterOperationOptions struct {
	Expand  *string
	Version *[]string
}

func DefaultDeploymentsListForClusterOperationOptions() DeploymentsListForClusterOperationOptions {
	return DeploymentsListForClusterOperationOptions{}
}

func (o DeploymentsListForClusterOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeploymentsListForClusterOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DeploymentsListForClusterOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Version != nil {
		out.Append("version", fmt.Sprintf("%v", *o.Version))
	}
	return &out
}

type DeploymentsListForClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DeploymentsListForClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DeploymentsListForCluster ...
func (c AppPlatformClient) DeploymentsListForCluster(ctx context.Context, id commonids.SpringCloudServiceId, options DeploymentsListForClusterOperationOptions) (result DeploymentsListForClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &DeploymentsListForClusterCustomPager{},
		Path:          fmt.Sprintf("%s/deployments", id.ID()),
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

// DeploymentsListForClusterComplete retrieves all the results into a single object
func (c AppPlatformClient) DeploymentsListForClusterComplete(ctx context.Context, id commonids.SpringCloudServiceId, options DeploymentsListForClusterOperationOptions) (DeploymentsListForClusterCompleteResult, error) {
	return c.DeploymentsListForClusterCompleteMatchingPredicate(ctx, id, options, DeploymentResourceOperationPredicate{})
}

// DeploymentsListForClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) DeploymentsListForClusterCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, options DeploymentsListForClusterOperationOptions, predicate DeploymentResourceOperationPredicate) (result DeploymentsListForClusterCompleteResult, err error) {
	items := make([]DeploymentResource, 0)

	resp, err := c.DeploymentsListForCluster(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = DeploymentsListForClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
