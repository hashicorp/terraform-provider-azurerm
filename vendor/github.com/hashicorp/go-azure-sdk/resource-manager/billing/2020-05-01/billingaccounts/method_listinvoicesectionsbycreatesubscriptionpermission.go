package billingaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListInvoiceSectionsByCreateSubscriptionPermissionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]InvoiceSectionWithCreateSubPermission
}

type ListInvoiceSectionsByCreateSubscriptionPermissionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []InvoiceSectionWithCreateSubPermission
}

// ListInvoiceSectionsByCreateSubscriptionPermission ...
func (c BillingAccountsClient) ListInvoiceSectionsByCreateSubscriptionPermission(ctx context.Context, id BillingAccountId) (result ListInvoiceSectionsByCreateSubscriptionPermissionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listInvoiceSectionsWithCreateSubscriptionPermission", id.ID()),
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
		Values *[]InvoiceSectionWithCreateSubPermission `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInvoiceSectionsByCreateSubscriptionPermissionComplete retrieves all the results into a single object
func (c BillingAccountsClient) ListInvoiceSectionsByCreateSubscriptionPermissionComplete(ctx context.Context, id BillingAccountId) (ListInvoiceSectionsByCreateSubscriptionPermissionCompleteResult, error) {
	return c.ListInvoiceSectionsByCreateSubscriptionPermissionCompleteMatchingPredicate(ctx, id, InvoiceSectionWithCreateSubPermissionOperationPredicate{})
}

// ListInvoiceSectionsByCreateSubscriptionPermissionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BillingAccountsClient) ListInvoiceSectionsByCreateSubscriptionPermissionCompleteMatchingPredicate(ctx context.Context, id BillingAccountId, predicate InvoiceSectionWithCreateSubPermissionOperationPredicate) (result ListInvoiceSectionsByCreateSubscriptionPermissionCompleteResult, err error) {
	items := make([]InvoiceSectionWithCreateSubPermission, 0)

	resp, err := c.ListInvoiceSectionsByCreateSubscriptionPermission(ctx, id)
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

	result = ListInvoiceSectionsByCreateSubscriptionPermissionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
