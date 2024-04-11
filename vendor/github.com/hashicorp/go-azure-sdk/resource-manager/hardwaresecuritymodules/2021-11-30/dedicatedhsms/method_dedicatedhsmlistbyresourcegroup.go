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

type DedicatedHsmListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DedicatedHsm
}

type DedicatedHsmListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DedicatedHsm
}

type DedicatedHsmListByResourceGroupOperationOptions struct {
	Top *int64
}

func DefaultDedicatedHsmListByResourceGroupOperationOptions() DedicatedHsmListByResourceGroupOperationOptions {
	return DedicatedHsmListByResourceGroupOperationOptions{}
}

func (o DedicatedHsmListByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DedicatedHsmListByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DedicatedHsmListByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// DedicatedHsmListByResourceGroup ...
func (c DedicatedHsmsClient) DedicatedHsmListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions) (result DedicatedHsmListByResourceGroupOperationResponse, err error) {
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

// DedicatedHsmListByResourceGroupComplete retrieves all the results into a single object
func (c DedicatedHsmsClient) DedicatedHsmListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions) (DedicatedHsmListByResourceGroupCompleteResult, error) {
	return c.DedicatedHsmListByResourceGroupCompleteMatchingPredicate(ctx, id, options, DedicatedHsmOperationPredicate{})
}

// DedicatedHsmListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DedicatedHsmsClient) DedicatedHsmListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions, predicate DedicatedHsmOperationPredicate) (result DedicatedHsmListByResourceGroupCompleteResult, err error) {
	items := make([]DedicatedHsm, 0)

	resp, err := c.DedicatedHsmListByResourceGroup(ctx, id, options)
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

	result = DedicatedHsmListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
