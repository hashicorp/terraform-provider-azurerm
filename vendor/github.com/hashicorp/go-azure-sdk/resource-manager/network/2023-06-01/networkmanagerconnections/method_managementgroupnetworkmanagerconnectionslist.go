package networkmanagerconnections

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

type ManagementGroupNetworkManagerConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkManagerConnection
}

type ManagementGroupNetworkManagerConnectionsListCompleteResult struct {
	Items []NetworkManagerConnection
}

type ManagementGroupNetworkManagerConnectionsListOperationOptions struct {
	Top *int64
}

func DefaultManagementGroupNetworkManagerConnectionsListOperationOptions() ManagementGroupNetworkManagerConnectionsListOperationOptions {
	return ManagementGroupNetworkManagerConnectionsListOperationOptions{}
}

func (o ManagementGroupNetworkManagerConnectionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ManagementGroupNetworkManagerConnectionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ManagementGroupNetworkManagerConnectionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ManagementGroupNetworkManagerConnectionsList ...
func (c NetworkManagerConnectionsClient) ManagementGroupNetworkManagerConnectionsList(ctx context.Context, id commonids.ManagementGroupId, options ManagementGroupNetworkManagerConnectionsListOperationOptions) (result ManagementGroupNetworkManagerConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Network/networkManagerConnections", id.ID()),
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
		Values *[]NetworkManagerConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ManagementGroupNetworkManagerConnectionsListComplete retrieves all the results into a single object
func (c NetworkManagerConnectionsClient) ManagementGroupNetworkManagerConnectionsListComplete(ctx context.Context, id commonids.ManagementGroupId, options ManagementGroupNetworkManagerConnectionsListOperationOptions) (ManagementGroupNetworkManagerConnectionsListCompleteResult, error) {
	return c.ManagementGroupNetworkManagerConnectionsListCompleteMatchingPredicate(ctx, id, options, NetworkManagerConnectionOperationPredicate{})
}

// ManagementGroupNetworkManagerConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkManagerConnectionsClient) ManagementGroupNetworkManagerConnectionsListCompleteMatchingPredicate(ctx context.Context, id commonids.ManagementGroupId, options ManagementGroupNetworkManagerConnectionsListOperationOptions, predicate NetworkManagerConnectionOperationPredicate) (result ManagementGroupNetworkManagerConnectionsListCompleteResult, err error) {
	items := make([]NetworkManagerConnection, 0)

	resp, err := c.ManagementGroupNetworkManagerConnectionsList(ctx, id, options)
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

	result = ManagementGroupNetworkManagerConnectionsListCompleteResult{
		Items: items,
	}
	return
}
