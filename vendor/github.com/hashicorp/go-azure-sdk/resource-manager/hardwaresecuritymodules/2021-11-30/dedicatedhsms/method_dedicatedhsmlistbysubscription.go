package dedicatedhsms

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

type DedicatedHsmListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DedicatedHsm
}

type DedicatedHsmListBySubscriptionCompleteResult struct {
	Items []DedicatedHsm
}

type DedicatedHsmListBySubscriptionOperationOptions struct {
	Top *int64
}

func DefaultDedicatedHsmListBySubscriptionOperationOptions() DedicatedHsmListBySubscriptionOperationOptions {
	return DedicatedHsmListBySubscriptionOperationOptions{}
}

func (o DedicatedHsmListBySubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DedicatedHsmListBySubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DedicatedHsmListBySubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// DedicatedHsmListBySubscription ...
func (c DedicatedHsmsClient) DedicatedHsmListBySubscription(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions) (result DedicatedHsmListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.HardwareSecurityModules/dedicatedHSMs", id.ID()),
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
		Values *[]DedicatedHsm `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DedicatedHsmListBySubscriptionComplete retrieves all the results into a single object
func (c DedicatedHsmsClient) DedicatedHsmListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions) (DedicatedHsmListBySubscriptionCompleteResult, error) {
	return c.DedicatedHsmListBySubscriptionCompleteMatchingPredicate(ctx, id, options, DedicatedHsmOperationPredicate{})
}

// DedicatedHsmListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DedicatedHsmsClient) DedicatedHsmListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions, predicate DedicatedHsmOperationPredicate) (result DedicatedHsmListBySubscriptionCompleteResult, err error) {
	items := make([]DedicatedHsm, 0)

	resp, err := c.DedicatedHsmListBySubscription(ctx, id, options)
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

	result = DedicatedHsmListBySubscriptionCompleteResult{
		Items: items,
	}
	return
}
