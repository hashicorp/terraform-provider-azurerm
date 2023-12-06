package dscconfiguration

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAutomationAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DscConfiguration
}

type ListByAutomationAccountCompleteResult struct {
	Items []DscConfiguration
}

type ListByAutomationAccountOperationOptions struct {
	Filter      *string
	Inlinecount *string
	Skip        *int64
	Top         *int64
}

func DefaultListByAutomationAccountOperationOptions() ListByAutomationAccountOperationOptions {
	return ListByAutomationAccountOperationOptions{}
}

func (o ListByAutomationAccountOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByAutomationAccountOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByAutomationAccountOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Inlinecount != nil {
		out.Append("$inlinecount", fmt.Sprintf("%v", *o.Inlinecount))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByAutomationAccount ...
func (c DscConfigurationClient) ListByAutomationAccount(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions) (result ListByAutomationAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/configurations", id.ID()),
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
		Values *[]DscConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAutomationAccountComplete retrieves all the results into a single object
func (c DscConfigurationClient) ListByAutomationAccountComplete(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions) (ListByAutomationAccountCompleteResult, error) {
	return c.ListByAutomationAccountCompleteMatchingPredicate(ctx, id, options, DscConfigurationOperationPredicate{})
}

// ListByAutomationAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DscConfigurationClient) ListByAutomationAccountCompleteMatchingPredicate(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions, predicate DscConfigurationOperationPredicate) (result ListByAutomationAccountCompleteResult, err error) {
	items := make([]DscConfiguration, 0)

	resp, err := c.ListByAutomationAccount(ctx, id, options)
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

	result = ListByAutomationAccountCompleteResult{
		Items: items,
	}
	return
}
