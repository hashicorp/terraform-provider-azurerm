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

type RemediationsListForManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Remediation
}

type RemediationsListForManagementGroupCompleteResult struct {
	Items []Remediation
}

type RemediationsListForManagementGroupOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRemediationsListForManagementGroupOperationOptions() RemediationsListForManagementGroupOperationOptions {
	return RemediationsListForManagementGroupOperationOptions{}
}

func (o RemediationsListForManagementGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RemediationsListForManagementGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RemediationsListForManagementGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// RemediationsListForManagementGroup ...
func (c RemediationsClient) RemediationsListForManagementGroup(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions) (result RemediationsListForManagementGroupOperationResponse, err error) {
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

// RemediationsListForManagementGroupComplete retrieves all the results into a single object
func (c RemediationsClient) RemediationsListForManagementGroupComplete(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions) (RemediationsListForManagementGroupCompleteResult, error) {
	return c.RemediationsListForManagementGroupCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// RemediationsListForManagementGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RemediationsClient) RemediationsListForManagementGroupCompleteMatchingPredicate(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions, predicate RemediationOperationPredicate) (result RemediationsListForManagementGroupCompleteResult, err error) {
	items := make([]Remediation, 0)

	resp, err := c.RemediationsListForManagementGroup(ctx, id, options)
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

	result = RemediationsListForManagementGroupCompleteResult{
		Items: items,
	}
	return
}
