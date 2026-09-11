package virtualmachines

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

type ListAllOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualMachine
}

type ListAllCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualMachine
}

type ListAllOperationOptions struct {
	Expand     *ExpandTypesForListVMs
	Filter     *string
	StatusOnly *string
}

func DefaultListAllOperationOptions() ListAllOperationOptions {
	return ListAllOperationOptions{}
}

func (o ListAllOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAllOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAllOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.StatusOnly != nil {
		out.Append("statusOnly", fmt.Sprintf("%v", *o.StatusOnly))
	}
	return &out
}

type ListAllCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAllCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAll ...
func (c VirtualMachinesClient) ListAll(ctx context.Context, id commonids.SubscriptionId, options ListAllOperationOptions) (result ListAllOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAllCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Compute/virtualMachines", id.ID()),
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
		Values *[]VirtualMachine `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAllComplete retrieves all the results into a single object
func (c VirtualMachinesClient) ListAllComplete(ctx context.Context, id commonids.SubscriptionId, options ListAllOperationOptions) (ListAllCompleteResult, error) {
	return c.ListAllCompleteMatchingPredicate(ctx, id, options, VirtualMachineOperationPredicate{})
}

// ListAllCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualMachinesClient) ListAllCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListAllOperationOptions, predicate VirtualMachineOperationPredicate) (result ListAllCompleteResult, err error) {
	items := make([]VirtualMachine, 0)

	resp, err := c.ListAll(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = ListAllCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
