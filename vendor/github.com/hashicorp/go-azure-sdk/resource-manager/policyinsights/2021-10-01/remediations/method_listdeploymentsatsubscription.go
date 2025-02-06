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

type ListDeploymentsAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemediationDeployment
}

type ListDeploymentsAtSubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RemediationDeployment
}

type ListDeploymentsAtSubscriptionOperationOptions struct {
	Top *int64
}

func DefaultListDeploymentsAtSubscriptionOperationOptions() ListDeploymentsAtSubscriptionOperationOptions {
	return ListDeploymentsAtSubscriptionOperationOptions{}
}

func (o ListDeploymentsAtSubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListDeploymentsAtSubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListDeploymentsAtSubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListDeploymentsAtSubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListDeploymentsAtSubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListDeploymentsAtSubscription ...
func (c RemediationsClient) ListDeploymentsAtSubscription(ctx context.Context, id RemediationId, options ListDeploymentsAtSubscriptionOperationOptions) (result ListDeploymentsAtSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListDeploymentsAtSubscriptionCustomPager{},
		Path:          fmt.Sprintf("%s/listDeployments", id.ID()),
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

// ListDeploymentsAtSubscriptionComplete retrieves all the results into a single object
func (c RemediationsClient) ListDeploymentsAtSubscriptionComplete(ctx context.Context, id RemediationId, options ListDeploymentsAtSubscriptionOperationOptions) (ListDeploymentsAtSubscriptionCompleteResult, error) {
	return c.ListDeploymentsAtSubscriptionCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// ListDeploymentsAtSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) ListDeploymentsAtSubscriptionCompleteMatchingPredicate(ctx context.Context, id RemediationId, options ListDeploymentsAtSubscriptionOperationOptions, predicate RemediationDeploymentOperationPredicate) (result ListDeploymentsAtSubscriptionCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	resp, err := c.ListDeploymentsAtSubscription(ctx, id, options)
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

	result = ListDeploymentsAtSubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
