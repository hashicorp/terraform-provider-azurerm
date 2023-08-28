package confidentialledger

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

type LedgerListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConfidentialLedger
}

type LedgerListByResourceGroupCompleteResult struct {
	Items []ConfidentialLedger
}

type LedgerListByResourceGroupOperationOptions struct {
	Filter *string
}

func DefaultLedgerListByResourceGroupOperationOptions() LedgerListByResourceGroupOperationOptions {
	return LedgerListByResourceGroupOperationOptions{}
}

func (o LedgerListByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LedgerListByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o LedgerListByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// LedgerListByResourceGroup ...
func (c ConfidentialLedgerClient) LedgerListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options LedgerListByResourceGroupOperationOptions) (result LedgerListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.ConfidentialLedger/ledgers", id.ID()),
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
		Values *[]ConfidentialLedger `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LedgerListByResourceGroupComplete retrieves all the results into a single object
func (c ConfidentialLedgerClient) LedgerListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options LedgerListByResourceGroupOperationOptions) (LedgerListByResourceGroupCompleteResult, error) {
	return c.LedgerListByResourceGroupCompleteMatchingPredicate(ctx, id, options, ConfidentialLedgerOperationPredicate{})
}

// LedgerListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConfidentialLedgerClient) LedgerListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options LedgerListByResourceGroupOperationOptions, predicate ConfidentialLedgerOperationPredicate) (result LedgerListByResourceGroupCompleteResult, err error) {
	items := make([]ConfidentialLedger, 0)

	resp, err := c.LedgerListByResourceGroup(ctx, id, options)
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

	result = LedgerListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
