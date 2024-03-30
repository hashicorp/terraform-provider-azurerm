package enrollmentaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByBillingAccountNameOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EnrollmentAccount
}

type ListByBillingAccountNameCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EnrollmentAccount
}

type ListByBillingAccountNameOperationOptions struct {
	Expand *string
	Filter *string
}

func DefaultListByBillingAccountNameOperationOptions() ListByBillingAccountNameOperationOptions {
	return ListByBillingAccountNameOperationOptions{}
}

func (o ListByBillingAccountNameOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByBillingAccountNameOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByBillingAccountNameOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ListByBillingAccountName ...
func (c EnrollmentAccountsClient) ListByBillingAccountName(ctx context.Context, id BillingAccountId, options ListByBillingAccountNameOperationOptions) (result ListByBillingAccountNameOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/enrollmentAccounts", id.ID()),
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
		Values *[]EnrollmentAccount `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByBillingAccountNameComplete retrieves all the results into a single object
func (c EnrollmentAccountsClient) ListByBillingAccountNameComplete(ctx context.Context, id BillingAccountId, options ListByBillingAccountNameOperationOptions) (ListByBillingAccountNameCompleteResult, error) {
	return c.ListByBillingAccountNameCompleteMatchingPredicate(ctx, id, options, EnrollmentAccountOperationPredicate{})
}

// ListByBillingAccountNameCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnrollmentAccountsClient) ListByBillingAccountNameCompleteMatchingPredicate(ctx context.Context, id BillingAccountId, options ListByBillingAccountNameOperationOptions, predicate EnrollmentAccountOperationPredicate) (result ListByBillingAccountNameCompleteResult, err error) {
	items := make([]EnrollmentAccount, 0)

	resp, err := c.ListByBillingAccountName(ctx, id, options)
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

	result = ListByBillingAccountNameCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
