package virtualnetworkappliancecapabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualNetworkApplianceCapability
}

type ListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualNetworkApplianceCapability
}

type ListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// List ...
func (c VirtualNetworkApplianceCapabilitiesClient) List(ctx context.Context, id VirtualNetworkApplianceId) (result ListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListCustomPager{},
		Path:       fmt.Sprintf("%s/capabilities", id.ID()),
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
		Values *[]json.RawMessage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	temp := make([]VirtualNetworkApplianceCapability, 0)
	if values.Values != nil {
		for i, v := range *values.Values {
			val, err := UnmarshalVirtualNetworkApplianceCapabilityImplementation(v)
			if err != nil {
				err = fmt.Errorf("unmarshalling item %d for VirtualNetworkApplianceCapability (%q): %+v", i, v, err)
				return result, err
			}
			temp = append(temp, val)
		}
	}
	result.Model = &temp

	return
}

// ListComplete retrieves all the results into a single object
func (c VirtualNetworkApplianceCapabilitiesClient) ListComplete(ctx context.Context, id VirtualNetworkApplianceId) (ListCompleteResult, error) {
	return c.ListCompleteMatchingPredicate(ctx, id, VirtualNetworkApplianceCapabilityOperationPredicate{})
}

// ListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualNetworkApplianceCapabilitiesClient) ListCompleteMatchingPredicate(ctx context.Context, id VirtualNetworkApplianceId, predicate VirtualNetworkApplianceCapabilityOperationPredicate) (result ListCompleteResult, err error) {
	items := make([]VirtualNetworkApplianceCapability, 0)

	resp, err := c.List(ctx, id)
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

	result = ListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
