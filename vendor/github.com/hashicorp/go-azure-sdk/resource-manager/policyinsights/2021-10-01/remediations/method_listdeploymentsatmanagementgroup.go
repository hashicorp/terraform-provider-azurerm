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

type ListDeploymentsAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemediationDeployment
}

type ListDeploymentsAtManagementGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RemediationDeployment
}

type ListDeploymentsAtManagementGroupOperationOptions struct {
	Top *int64
}

func DefaultListDeploymentsAtManagementGroupOperationOptions() ListDeploymentsAtManagementGroupOperationOptions {
	return ListDeploymentsAtManagementGroupOperationOptions{}
}

func (o ListDeploymentsAtManagementGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListDeploymentsAtManagementGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListDeploymentsAtManagementGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListDeploymentsAtManagementGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListDeploymentsAtManagementGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListDeploymentsAtManagementGroup ...
func (c RemediationsClient) ListDeploymentsAtManagementGroup(ctx context.Context, id Providers2RemediationId, options ListDeploymentsAtManagementGroupOperationOptions) (result ListDeploymentsAtManagementGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListDeploymentsAtManagementGroupCustomPager{},
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

// ListDeploymentsAtManagementGroupComplete retrieves all the results into a single object
func (c RemediationsClient) ListDeploymentsAtManagementGroupComplete(ctx context.Context, id Providers2RemediationId, options ListDeploymentsAtManagementGroupOperationOptions) (ListDeploymentsAtManagementGroupCompleteResult, error) {
	return c.ListDeploymentsAtManagementGroupCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// ListDeploymentsAtManagementGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) ListDeploymentsAtManagementGroupCompleteMatchingPredicate(ctx context.Context, id Providers2RemediationId, options ListDeploymentsAtManagementGroupOperationOptions, predicate RemediationDeploymentOperationPredicate) (result ListDeploymentsAtManagementGroupCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	resp, err := c.ListDeploymentsAtManagementGroup(ctx, id, options)
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

	result = ListDeploymentsAtManagementGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
