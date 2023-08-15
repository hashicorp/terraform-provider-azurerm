package remediations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsListDeploymentsAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemediationDeployment
}

type RemediationsListDeploymentsAtSubscriptionCompleteResult struct {
	Items []RemediationDeployment
}

type RemediationsListDeploymentsAtSubscriptionOperationOptions struct {
	Top *int64
}

func DefaultRemediationsListDeploymentsAtSubscriptionOperationOptions() RemediationsListDeploymentsAtSubscriptionOperationOptions {
	return RemediationsListDeploymentsAtSubscriptionOperationOptions{}
}

func (o RemediationsListDeploymentsAtSubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RemediationsListDeploymentsAtSubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RemediationsListDeploymentsAtSubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RemediationsListDeploymentsAtSubscription ...
func (c RemediationsClient) RemediationsListDeploymentsAtSubscription(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions) (result RemediationsListDeploymentsAtSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/listDeployments", id.ID()),
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
		Values *[]RemediationDeployment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RemediationsListDeploymentsAtSubscriptionComplete retrieves all the results into a single object
func (c RemediationsClient) RemediationsListDeploymentsAtSubscriptionComplete(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions) (RemediationsListDeploymentsAtSubscriptionCompleteResult, error) {
	return c.RemediationsListDeploymentsAtSubscriptionCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// RemediationsListDeploymentsAtSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) RemediationsListDeploymentsAtSubscriptionCompleteMatchingPredicate(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions, predicate RemediationDeploymentOperationPredicate) (result RemediationsListDeploymentsAtSubscriptionCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	resp, err := c.RemediationsListDeploymentsAtSubscription(ctx, id, options)
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

	result = RemediationsListDeploymentsAtSubscriptionCompleteResult{
		Items: items,
	}
	return
}
