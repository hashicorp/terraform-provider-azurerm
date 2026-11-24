package dscpconfigurations

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

type DscpConfigurationListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DscpConfiguration
}

type DscpConfigurationListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DscpConfiguration
}

type DscpConfigurationListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DscpConfigurationListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DscpConfigurationList ...
func (c DscpConfigurationsClient) DscpConfigurationList(ctx context.Context, id commonids.ResourceGroupId) (result DscpConfigurationListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DscpConfigurationListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/dscpConfigurations", id.ID()),
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
		Values *[]DscpConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DscpConfigurationListComplete retrieves all the results into a single object
func (c DscpConfigurationsClient) DscpConfigurationListComplete(ctx context.Context, id commonids.ResourceGroupId) (DscpConfigurationListCompleteResult, error) {
	return c.DscpConfigurationListCompleteMatchingPredicate(ctx, id, DscpConfigurationOperationPredicate{})
}

// DscpConfigurationListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DscpConfigurationsClient) DscpConfigurationListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate DscpConfigurationOperationPredicate) (result DscpConfigurationListCompleteResult, err error) {
	items := make([]DscpConfiguration, 0)

	resp, err := c.DscpConfigurationList(ctx, id)
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

	result = DscpConfigurationListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
