package roleassignments

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
	Model        *[]RoleAssignment
}

type ListForSubscriptionCompleteResult struct {
	Items []RoleAssignment
}

type ListForSubscriptionOperationOptions struct {
	Filter   *string
	TenantId *string
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
	if o.TenantId != nil {
		out.Append("tenantId", fmt.Sprintf("%v", *o.TenantId))
	}
	return &out
}

// ListForSubscription ...
func (c RoleAssignmentsClient) ListForSubscription(ctx context.Context, id commonids.SubscriptionId, options ListForSubscriptionOperationOptions) (result ListForSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments", id.ID()),
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
		Values *[]RoleAssignment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListForSubscriptionComplete retrieves all the results into a single object
func (c RoleAssignmentsClient) ListForSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options ListForSubscriptionOperationOptions) (ListForSubscriptionCompleteResult, error) {
	return c.ListForSubscriptionCompleteMatchingPredicate(ctx, id, options, RoleAssignmentOperationPredicate{})
}

// ListForSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RoleAssignmentsClient) ListForSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListForSubscriptionOperationOptions, predicate RoleAssignmentOperationPredicate) (result ListForSubscriptionCompleteResult, err error) {
	items := make([]RoleAssignment, 0)

	resp, err := c.ListForSubscription(ctx, id, options)
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

	result = ListForSubscriptionCompleteResult{
		Items: items,
	}
	return
}
