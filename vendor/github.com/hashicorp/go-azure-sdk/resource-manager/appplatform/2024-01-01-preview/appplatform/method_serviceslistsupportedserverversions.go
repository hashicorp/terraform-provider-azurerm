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

type ServicesListSupportedServerVersionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SupportedServerVersion
}

type ServicesListSupportedServerVersionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SupportedServerVersion
}

type ServicesListSupportedServerVersionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServicesListSupportedServerVersionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServicesListSupportedServerVersions ...
func (c AppPlatformClient) ServicesListSupportedServerVersions(ctx context.Context, id commonids.SpringCloudServiceId) (result ServicesListSupportedServerVersionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ServicesListSupportedServerVersionsCustomPager{},
		Path:       fmt.Sprintf("%s/supportedServerVersions", id.ID()),
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
		Values *[]SupportedServerVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServicesListSupportedServerVersionsComplete retrieves all the results into a single object
func (c AppPlatformClient) ServicesListSupportedServerVersionsComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ServicesListSupportedServerVersionsCompleteResult, error) {
	return c.ServicesListSupportedServerVersionsCompleteMatchingPredicate(ctx, id, SupportedServerVersionOperationPredicate{})
}

// ServicesListSupportedServerVersionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ServicesListSupportedServerVersionsCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate SupportedServerVersionOperationPredicate) (result ServicesListSupportedServerVersionsCompleteResult, err error) {
	items := make([]SupportedServerVersion, 0)

	resp, err := c.ServicesListSupportedServerVersions(ctx, id)
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

	result = ServicesListSupportedServerVersionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
