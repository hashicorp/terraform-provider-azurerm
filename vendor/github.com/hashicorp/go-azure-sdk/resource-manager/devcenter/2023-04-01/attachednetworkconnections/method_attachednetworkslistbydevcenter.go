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

type AttachedNetworksListByDevCenterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AttachedNetworkConnection
}

type AttachedNetworksListByDevCenterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AttachedNetworkConnection
}

type AttachedNetworksListByDevCenterOperationOptions struct {
	Top *int64
}

func DefaultAttachedNetworksListByDevCenterOperationOptions() AttachedNetworksListByDevCenterOperationOptions {
	return AttachedNetworksListByDevCenterOperationOptions{}
}

func (o AttachedNetworksListByDevCenterOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AttachedNetworksListByDevCenterOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AttachedNetworksListByDevCenterOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type AttachedNetworksListByDevCenterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AttachedNetworksListByDevCenterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AttachedNetworksListByDevCenter ...
func (c AttachedNetworkConnectionsClient) AttachedNetworksListByDevCenter(ctx context.Context, id DevCenterId, options AttachedNetworksListByDevCenterOperationOptions) (result AttachedNetworksListByDevCenterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &AttachedNetworksListByDevCenterCustomPager{},
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

// AttachedNetworksListByDevCenterComplete retrieves all the results into a single object
func (c AttachedNetworkConnectionsClient) AttachedNetworksListByDevCenterComplete(ctx context.Context, id DevCenterId, options AttachedNetworksListByDevCenterOperationOptions) (AttachedNetworksListByDevCenterCompleteResult, error) {
	return c.AttachedNetworksListByDevCenterCompleteMatchingPredicate(ctx, id, options, AttachedNetworkConnectionOperationPredicate{})
}

// AttachedNetworksListByDevCenterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AttachedNetworkConnectionsClient) AttachedNetworksListByDevCenterCompleteMatchingPredicate(ctx context.Context, id DevCenterId, options AttachedNetworksListByDevCenterOperationOptions, predicate AttachedNetworkConnectionOperationPredicate) (result AttachedNetworksListByDevCenterCompleteResult, err error) {
	items := make([]AttachedNetworkConnection, 0)

	resp, err := c.AttachedNetworksListByDevCenter(ctx, id, options)
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

	result = AttachedNetworksListByDevCenterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
