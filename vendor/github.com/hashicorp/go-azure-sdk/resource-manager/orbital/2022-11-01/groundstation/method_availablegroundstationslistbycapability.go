package groundstation

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

type AvailableGroundStationsListByCapabilityOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailableGroundStation
}

type AvailableGroundStationsListByCapabilityCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AvailableGroundStation
}

type AvailableGroundStationsListByCapabilityOperationOptions struct {
	Capability *CapabilityParameter
}

func DefaultAvailableGroundStationsListByCapabilityOperationOptions() AvailableGroundStationsListByCapabilityOperationOptions {
	return AvailableGroundStationsListByCapabilityOperationOptions{}
}

func (o AvailableGroundStationsListByCapabilityOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AvailableGroundStationsListByCapabilityOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AvailableGroundStationsListByCapabilityOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Capability != nil {
		out.Append("capability", fmt.Sprintf("%v", *o.Capability))
	}
	return &out
}

type AvailableGroundStationsListByCapabilityCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AvailableGroundStationsListByCapabilityCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AvailableGroundStationsListByCapability ...
func (c GroundStationClient) AvailableGroundStationsListByCapability(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions) (result AvailableGroundStationsListByCapabilityOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &AvailableGroundStationsListByCapabilityCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Orbital/availableGroundStations", id.ID()),
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
		Values *[]AvailableGroundStation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AvailableGroundStationsListByCapabilityComplete retrieves all the results into a single object
func (c GroundStationClient) AvailableGroundStationsListByCapabilityComplete(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions) (AvailableGroundStationsListByCapabilityCompleteResult, error) {
	return c.AvailableGroundStationsListByCapabilityCompleteMatchingPredicate(ctx, id, options, AvailableGroundStationOperationPredicate{})
}

// AvailableGroundStationsListByCapabilityCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GroundStationClient) AvailableGroundStationsListByCapabilityCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions, predicate AvailableGroundStationOperationPredicate) (result AvailableGroundStationsListByCapabilityCompleteResult, err error) {
	items := make([]AvailableGroundStation, 0)

	resp, err := c.AvailableGroundStationsListByCapability(ctx, id, options)
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

	result = AvailableGroundStationsListByCapabilityCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
