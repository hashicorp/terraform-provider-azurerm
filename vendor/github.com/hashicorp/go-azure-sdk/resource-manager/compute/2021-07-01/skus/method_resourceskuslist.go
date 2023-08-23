package skus

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

type ResourceSkusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceSku
}

type ResourceSkusListCompleteResult struct {
	Items []ResourceSku
}

type ResourceSkusListOperationOptions struct {
	Filter                   *string
	IncludeExtendedLocations *string
}

func DefaultResourceSkusListOperationOptions() ResourceSkusListOperationOptions {
	return ResourceSkusListOperationOptions{}
}

func (o ResourceSkusListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ResourceSkusListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ResourceSkusListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IncludeExtendedLocations != nil {
		out.Append("includeExtendedLocations", fmt.Sprintf("%v", *o.IncludeExtendedLocations))
	}
	return &out
}

// ResourceSkusList ...
func (c SkusClient) ResourceSkusList(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions) (result ResourceSkusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Compute/skus", id.ID()),
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
		Values *[]ResourceSku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ResourceSkusListComplete retrieves all the results into a single object
func (c SkusClient) ResourceSkusListComplete(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions) (ResourceSkusListCompleteResult, error) {
	return c.ResourceSkusListCompleteMatchingPredicate(ctx, id, options, ResourceSkuOperationPredicate{})
}

// ResourceSkusListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SkusClient) ResourceSkusListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions, predicate ResourceSkuOperationPredicate) (result ResourceSkusListCompleteResult, err error) {
	items := make([]ResourceSku, 0)

	resp, err := c.ResourceSkusList(ctx, id, options)
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

	result = ResourceSkusListCompleteResult{
		Items: items,
	}
	return
}
