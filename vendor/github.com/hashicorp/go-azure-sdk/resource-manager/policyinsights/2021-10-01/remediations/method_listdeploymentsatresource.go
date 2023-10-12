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

type ListDeploymentsAtResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemediationDeployment
}

type ListDeploymentsAtResourceCompleteResult struct {
	Items []RemediationDeployment
}

type ListDeploymentsAtResourceOperationOptions struct {
	Top *int64
}

func DefaultListDeploymentsAtResourceOperationOptions() ListDeploymentsAtResourceOperationOptions {
	return ListDeploymentsAtResourceOperationOptions{}
}

func (o ListDeploymentsAtResourceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListDeploymentsAtResourceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListDeploymentsAtResourceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListDeploymentsAtResource ...
func (c RemediationsClient) ListDeploymentsAtResource(ctx context.Context, id ScopedRemediationId, options ListDeploymentsAtResourceOperationOptions) (result ListDeploymentsAtResourceOperationResponse, err error) {
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

// ListDeploymentsAtResourceComplete retrieves all the results into a single object
func (c RemediationsClient) ListDeploymentsAtResourceComplete(ctx context.Context, id ScopedRemediationId, options ListDeploymentsAtResourceOperationOptions) (ListDeploymentsAtResourceCompleteResult, error) {
	return c.ListDeploymentsAtResourceCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// ListDeploymentsAtResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) ListDeploymentsAtResourceCompleteMatchingPredicate(ctx context.Context, id ScopedRemediationId, options ListDeploymentsAtResourceOperationOptions, predicate RemediationDeploymentOperationPredicate) (result ListDeploymentsAtResourceCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	resp, err := c.ListDeploymentsAtResource(ctx, id, options)
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

	result = ListDeploymentsAtResourceCompleteResult{
		Items: items,
	}
	return
}
