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

type RemediationsListForSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Remediation
}

type RemediationsListForSubscriptionCompleteResult struct {
	Items []Remediation
}

type RemediationsListForSubscriptionOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRemediationsListForSubscriptionOperationOptions() RemediationsListForSubscriptionOperationOptions {
	return RemediationsListForSubscriptionOperationOptions{}
}

func (o RemediationsListForSubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RemediationsListForSubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RemediationsListForSubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RemediationsListForSubscription ...
func (c RemediationsClient) RemediationsListForSubscription(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions) (result RemediationsListForSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.PolicyInsights/remediations", id.ID()),
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
		Values *[]Remediation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RemediationsListForSubscriptionComplete retrieves all the results into a single object
func (c RemediationsClient) RemediationsListForSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions) (RemediationsListForSubscriptionCompleteResult, error) {
	return c.RemediationsListForSubscriptionCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// RemediationsListForSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) RemediationsListForSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions, predicate RemediationOperationPredicate) (result RemediationsListForSubscriptionCompleteResult, err error) {
	items := make([]Remediation, 0)

	resp, err := c.RemediationsListForSubscription(ctx, id, options)
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

	result = RemediationsListForSubscriptionCompleteResult{
		Items: items,
	}
	return
}
