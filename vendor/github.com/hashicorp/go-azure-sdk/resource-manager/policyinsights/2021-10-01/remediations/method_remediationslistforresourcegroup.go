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

type RemediationsListForResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Remediation
}

type RemediationsListForResourceGroupCompleteResult struct {
	Items []Remediation
}

type RemediationsListForResourceGroupOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRemediationsListForResourceGroupOperationOptions() RemediationsListForResourceGroupOperationOptions {
	return RemediationsListForResourceGroupOperationOptions{}
}

func (o RemediationsListForResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RemediationsListForResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RemediationsListForResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RemediationsListForResourceGroup ...
func (c RemediationsClient) RemediationsListForResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options RemediationsListForResourceGroupOperationOptions) (result RemediationsListForResourceGroupOperationResponse, err error) {
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

// RemediationsListForResourceGroupComplete retrieves all the results into a single object
func (c RemediationsClient) RemediationsListForResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options RemediationsListForResourceGroupOperationOptions) (RemediationsListForResourceGroupCompleteResult, error) {
	return c.RemediationsListForResourceGroupCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// RemediationsListForResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) RemediationsListForResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options RemediationsListForResourceGroupOperationOptions, predicate RemediationOperationPredicate) (result RemediationsListForResourceGroupCompleteResult, err error) {
	items := make([]Remediation, 0)

	resp, err := c.RemediationsListForResourceGroup(ctx, id, options)
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

	result = RemediationsListForResourceGroupCompleteResult{
		Items: items,
	}
	return
}
