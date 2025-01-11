package attachednetworkconnections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachedNetworksListByProjectOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AttachedNetworkConnection
}

type AttachedNetworksListByProjectCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AttachedNetworkConnection
}

type AttachedNetworksListByProjectOperationOptions struct {
	Top *int64
}

func DefaultAttachedNetworksListByProjectOperationOptions() AttachedNetworksListByProjectOperationOptions {
	return AttachedNetworksListByProjectOperationOptions{}
}

func (o AttachedNetworksListByProjectOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AttachedNetworksListByProjectOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AttachedNetworksListByProjectOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type AttachedNetworksListByProjectCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AttachedNetworksListByProjectCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AttachedNetworksListByProject ...
func (c AttachedNetworkConnectionsClient) AttachedNetworksListByProject(ctx context.Context, id ProjectId, options AttachedNetworksListByProjectOperationOptions) (result AttachedNetworksListByProjectOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &AttachedNetworksListByProjectCustomPager{},
		Path:          fmt.Sprintf("%s/attachedNetworks", id.ID()),
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
		Values *[]AttachedNetworkConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AttachedNetworksListByProjectComplete retrieves all the results into a single object
func (c AttachedNetworkConnectionsClient) AttachedNetworksListByProjectComplete(ctx context.Context, id ProjectId, options AttachedNetworksListByProjectOperationOptions) (AttachedNetworksListByProjectCompleteResult, error) {
	return c.AttachedNetworksListByProjectCompleteMatchingPredicate(ctx, id, options, AttachedNetworkConnectionOperationPredicate{})
}

// AttachedNetworksListByProjectCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AttachedNetworkConnectionsClient) AttachedNetworksListByProjectCompleteMatchingPredicate(ctx context.Context, id ProjectId, options AttachedNetworksListByProjectOperationOptions, predicate AttachedNetworkConnectionOperationPredicate) (result AttachedNetworksListByProjectCompleteResult, err error) {
	items := make([]AttachedNetworkConnection, 0)

	resp, err := c.AttachedNetworksListByProject(ctx, id, options)
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

	result = AttachedNetworksListByProjectCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
