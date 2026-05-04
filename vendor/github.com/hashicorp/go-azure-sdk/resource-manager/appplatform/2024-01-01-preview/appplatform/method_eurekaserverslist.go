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

type EurekaServersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EurekaServerResource
}

type EurekaServersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EurekaServerResource
}

type EurekaServersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *EurekaServersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// EurekaServersList ...
func (c AppPlatformClient) EurekaServersList(ctx context.Context, id commonids.SpringCloudServiceId) (result EurekaServersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &EurekaServersListCustomPager{},
		Path:       fmt.Sprintf("%s/eurekaServers", id.ID()),
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
		Values *[]EurekaServerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// EurekaServersListComplete retrieves all the results into a single object
func (c AppPlatformClient) EurekaServersListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (EurekaServersListCompleteResult, error) {
	return c.EurekaServersListCompleteMatchingPredicate(ctx, id, EurekaServerResourceOperationPredicate{})
}

// EurekaServersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) EurekaServersListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate EurekaServerResourceOperationPredicate) (result EurekaServersListCompleteResult, err error) {
	items := make([]EurekaServerResource, 0)

	resp, err := c.EurekaServersList(ctx, id)
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

	result = EurekaServersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
