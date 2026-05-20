package resourcegroups

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

type ResourcesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GenericResourceExpanded
}

type ResourcesListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GenericResourceExpanded
}

type ResourcesListByResourceGroupOperationOptions struct {
	Expand *string
	Filter *string
	Top    *int64
}

func DefaultResourcesListByResourceGroupOperationOptions() ResourcesListByResourceGroupOperationOptions {
	return ResourcesListByResourceGroupOperationOptions{}
}

func (o ResourcesListByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ResourcesListByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ResourcesListByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ResourcesListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ResourcesListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ResourcesListByResourceGroup ...
func (c ResourceGroupsClient) ResourcesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ResourcesListByResourceGroupOperationOptions) (result ResourcesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ResourcesListByResourceGroupCustomPager{},
		Path:          fmt.Sprintf("%s/resources", id.ID()),
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
		Values *[]GenericResourceExpanded `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ResourcesListByResourceGroupComplete retrieves all the results into a single object
func (c ResourceGroupsClient) ResourcesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options ResourcesListByResourceGroupOperationOptions) (ResourcesListByResourceGroupCompleteResult, error) {
	return c.ResourcesListByResourceGroupCompleteMatchingPredicate(ctx, id, options, GenericResourceExpandedOperationPredicate{})
}

// ResourcesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGroupsClient) ResourcesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options ResourcesListByResourceGroupOperationOptions, predicate GenericResourceExpandedOperationPredicate) (result ResourcesListByResourceGroupCompleteResult, err error) {
	items := make([]GenericResourceExpanded, 0)

	resp, err := c.ResourcesListByResourceGroup(ctx, id, options)
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

	result = ResourcesListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
