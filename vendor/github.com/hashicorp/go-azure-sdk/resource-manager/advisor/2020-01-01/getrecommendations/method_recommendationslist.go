package getrecommendations

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

type RecommendationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceRecommendationBase
}

type RecommendationsListCompleteResult struct {
	Items []ResourceRecommendationBase
}

type RecommendationsListOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRecommendationsListOperationOptions() RecommendationsListOperationOptions {
	return RecommendationsListOperationOptions{}
}

func (o RecommendationsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RecommendationsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RecommendationsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RecommendationsList ...
func (c GetRecommendationsClient) RecommendationsList(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions) (result RecommendationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Advisor/recommendations", id.ID()),
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
		Values *[]ResourceRecommendationBase `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RecommendationsListComplete retrieves all the results into a single object
func (c GetRecommendationsClient) RecommendationsListComplete(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions) (RecommendationsListCompleteResult, error) {
	return c.RecommendationsListCompleteMatchingPredicate(ctx, id, options, ResourceRecommendationBaseOperationPredicate{})
}

// RecommendationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GetRecommendationsClient) RecommendationsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions, predicate ResourceRecommendationBaseOperationPredicate) (result RecommendationsListCompleteResult, err error) {
	items := make([]ResourceRecommendationBase, 0)

	resp, err := c.RecommendationsList(ctx, id, options)
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

	result = RecommendationsListCompleteResult{
		Items: items,
	}
	return
}
