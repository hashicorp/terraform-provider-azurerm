package appplatform

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

type ServicesListSupportedApmTypesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SupportedApmType
}

type ServicesListSupportedApmTypesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SupportedApmType
}

type ServicesListSupportedApmTypesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServicesListSupportedApmTypesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServicesListSupportedApmTypes ...
func (c AppPlatformClient) ServicesListSupportedApmTypes(ctx context.Context, id commonids.SpringCloudServiceId) (result ServicesListSupportedApmTypesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ServicesListSupportedApmTypesCustomPager{},
		Path:       fmt.Sprintf("%s/supportedApmTypes", id.ID()),
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
		Values *[]SupportedApmType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServicesListSupportedApmTypesComplete retrieves all the results into a single object
func (c AppPlatformClient) ServicesListSupportedApmTypesComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ServicesListSupportedApmTypesCompleteResult, error) {
	return c.ServicesListSupportedApmTypesCompleteMatchingPredicate(ctx, id, SupportedApmTypeOperationPredicate{})
}

// ServicesListSupportedApmTypesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ServicesListSupportedApmTypesCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate SupportedApmTypeOperationPredicate) (result ServicesListSupportedApmTypesCompleteResult, err error) {
	items := make([]SupportedApmType, 0)

	resp, err := c.ServicesListSupportedApmTypes(ctx, id)
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

	result = ServicesListSupportedApmTypesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
