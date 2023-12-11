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

type ListDeploymentsAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RemediationDeployment
}

type ListDeploymentsAtResourceGroupCompleteResult struct {
	Items []RemediationDeployment
}

type ListDeploymentsAtResourceGroupOperationOptions struct {
	Top *int64
}

func DefaultListDeploymentsAtResourceGroupOperationOptions() ListDeploymentsAtResourceGroupOperationOptions {
	return ListDeploymentsAtResourceGroupOperationOptions{}
}

func (o ListDeploymentsAtResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListDeploymentsAtResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListDeploymentsAtResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListDeploymentsAtResourceGroup ...
func (c RemediationsClient) ListDeploymentsAtResourceGroup(ctx context.Context, id ProviderRemediationId, options ListDeploymentsAtResourceGroupOperationOptions) (result ListDeploymentsAtResourceGroupOperationResponse, err error) {
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

// ListDeploymentsAtResourceGroupComplete retrieves all the results into a single object
func (c RemediationsClient) ListDeploymentsAtResourceGroupComplete(ctx context.Context, id ProviderRemediationId, options ListDeploymentsAtResourceGroupOperationOptions) (ListDeploymentsAtResourceGroupCompleteResult, error) {
	return c.ListDeploymentsAtResourceGroupCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// ListDeploymentsAtResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) ListDeploymentsAtResourceGroupCompleteMatchingPredicate(ctx context.Context, id ProviderRemediationId, options ListDeploymentsAtResourceGroupOperationOptions, predicate RemediationDeploymentOperationPredicate) (result ListDeploymentsAtResourceGroupCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	resp, err := c.ListDeploymentsAtResourceGroup(ctx, id, options)
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

	result = ListDeploymentsAtResourceGroupCompleteResult{
		Items: items,
	}
	return
}
