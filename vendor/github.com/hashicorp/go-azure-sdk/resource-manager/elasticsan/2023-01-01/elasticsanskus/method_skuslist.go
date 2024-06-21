package elasticsanskus

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

type SkusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SkuInformation
}

type SkusListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SkuInformation
}

type SkusListOperationOptions struct {
	Filter *string
}

func DefaultSkusListOperationOptions() SkusListOperationOptions {
	return SkusListOperationOptions{}
}

func (o SkusListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o SkusListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o SkusListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// SkusList ...
func (c ElasticSanSkusClient) SkusList(ctx context.Context, id commonids.SubscriptionId, options SkusListOperationOptions) (result SkusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.ElasticSan/skus", id.ID()),
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
		Values *[]SkuInformation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SkusListComplete retrieves all the results into a single object
func (c ElasticSanSkusClient) SkusListComplete(ctx context.Context, id commonids.SubscriptionId, options SkusListOperationOptions) (SkusListCompleteResult, error) {
	return c.SkusListCompleteMatchingPredicate(ctx, id, options, SkuInformationOperationPredicate{})
}

// SkusListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ElasticSanSkusClient) SkusListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options SkusListOperationOptions, predicate SkuInformationOperationPredicate) (result SkusListCompleteResult, err error) {
	items := make([]SkuInformation, 0)

	resp, err := c.SkusList(ctx, id, options)
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

	result = SkusListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
