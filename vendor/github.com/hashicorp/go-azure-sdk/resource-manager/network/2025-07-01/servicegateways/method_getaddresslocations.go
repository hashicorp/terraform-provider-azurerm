package servicegateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAddressLocationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServiceGatewayAddressLocationResponse
}

type GetAddressLocationsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServiceGatewayAddressLocationResponse
}

type GetAddressLocationsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetAddressLocationsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetAddressLocations ...
func (c ServiceGatewaysClient) GetAddressLocations(ctx context.Context, id ServiceGatewayId) (result GetAddressLocationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetAddressLocationsCustomPager{},
		Path:       fmt.Sprintf("%s/addressLocations", id.ID()),
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
		Values *[]ServiceGatewayAddressLocationResponse `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetAddressLocationsComplete retrieves all the results into a single object
func (c ServiceGatewaysClient) GetAddressLocationsComplete(ctx context.Context, id ServiceGatewayId) (GetAddressLocationsCompleteResult, error) {
	return c.GetAddressLocationsCompleteMatchingPredicate(ctx, id, ServiceGatewayAddressLocationResponseOperationPredicate{})
}

// GetAddressLocationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServiceGatewaysClient) GetAddressLocationsCompleteMatchingPredicate(ctx context.Context, id ServiceGatewayId, predicate ServiceGatewayAddressLocationResponseOperationPredicate) (result GetAddressLocationsCompleteResult, err error) {
	items := make([]ServiceGatewayAddressLocationResponse, 0)

	resp, err := c.GetAddressLocations(ctx, id)
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

	result = GetAddressLocationsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
