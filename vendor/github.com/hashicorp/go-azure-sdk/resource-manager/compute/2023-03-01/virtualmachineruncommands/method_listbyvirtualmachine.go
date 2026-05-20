package virtualmachineruncommands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByVirtualMachineOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualMachineRunCommand
}

type ListByVirtualMachineCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualMachineRunCommand
}

type ListByVirtualMachineOperationOptions struct {
	Expand *string
}

func DefaultListByVirtualMachineOperationOptions() ListByVirtualMachineOperationOptions {
	return ListByVirtualMachineOperationOptions{}
}

func (o ListByVirtualMachineOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByVirtualMachineOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByVirtualMachineOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

type ListByVirtualMachineCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByVirtualMachineCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByVirtualMachine ...
func (c VirtualMachineRunCommandsClient) ListByVirtualMachine(ctx context.Context, id VirtualMachineId, options ListByVirtualMachineOperationOptions) (result ListByVirtualMachineOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByVirtualMachineCustomPager{},
		Path:          fmt.Sprintf("%s/runCommands", id.ID()),
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
		Values *[]VirtualMachineRunCommand `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByVirtualMachineComplete retrieves all the results into a single object
func (c VirtualMachineRunCommandsClient) ListByVirtualMachineComplete(ctx context.Context, id VirtualMachineId, options ListByVirtualMachineOperationOptions) (ListByVirtualMachineCompleteResult, error) {
	return c.ListByVirtualMachineCompleteMatchingPredicate(ctx, id, options, VirtualMachineRunCommandOperationPredicate{})
}

// ListByVirtualMachineCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualMachineRunCommandsClient) ListByVirtualMachineCompleteMatchingPredicate(ctx context.Context, id VirtualMachineId, options ListByVirtualMachineOperationOptions, predicate VirtualMachineRunCommandOperationPredicate) (result ListByVirtualMachineCompleteResult, err error) {
	items := make([]VirtualMachineRunCommand, 0)

	resp, err := c.ListByVirtualMachine(ctx, id, options)
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

	result = ListByVirtualMachineCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
