package remediations

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

type ListForSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Remediation
}

type ListForSubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Remediation
}

type ListForSubscriptionOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListForSubscriptionOperationOptions() ListForSubscriptionOperationOptions {
	return ListForSubscriptionOperationOptions{}
}

func (o ListForSubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListForSubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListForSubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListForSubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListForSubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListForSubscription ...
func (c RemediationsClient) ListForSubscription(ctx context.Context, id commonids.SubscriptionId, options ListForSubscriptionOperationOptions) (result ListForSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListForSubscriptionCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.PolicyInsights/remediations", id.ID()),
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
		Values *[]Remediation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListForSubscriptionComplete retrieves all the results into a single object
func (c RemediationsClient) ListForSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options ListForSubscriptionOperationOptions) (ListForSubscriptionCompleteResult, error) {
	return c.ListForSubscriptionCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// ListForSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) ListForSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListForSubscriptionOperationOptions, predicate RemediationOperationPredicate) (result ListForSubscriptionCompleteResult, err error) {
	items := make([]Remediation, 0)

	resp, err := c.ListForSubscription(ctx, id, options)
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

	result = ListForSubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
