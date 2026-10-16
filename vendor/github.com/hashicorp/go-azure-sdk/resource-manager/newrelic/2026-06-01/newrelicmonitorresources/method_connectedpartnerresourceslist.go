package newrelicmonitorresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedPartnerResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConnectedPartnerResourcesListFormat
}

type ConnectedPartnerResourcesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConnectedPartnerResourcesListFormat
}

type ConnectedPartnerResourcesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ConnectedPartnerResourcesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ConnectedPartnerResourcesList ...
func (c NewRelicMonitorResourcesClient) ConnectedPartnerResourcesList(ctx context.Context, id MonitorId, input string) (result ConnectedPartnerResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ConnectedPartnerResourcesListCustomPager{},
		Path:       fmt.Sprintf("%s/listConnectedPartnerResources", id.ID()),
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
		Values *[]ConnectedPartnerResourcesListFormat `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ConnectedPartnerResourcesListComplete retrieves all the results into a single object
func (c NewRelicMonitorResourcesClient) ConnectedPartnerResourcesListComplete(ctx context.Context, id MonitorId, input string) (ConnectedPartnerResourcesListCompleteResult, error) {
	return c.ConnectedPartnerResourcesListCompleteMatchingPredicate(ctx, id, input, ConnectedPartnerResourcesListFormatOperationPredicate{})
}

// ConnectedPartnerResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NewRelicMonitorResourcesClient) ConnectedPartnerResourcesListCompleteMatchingPredicate(ctx context.Context, id MonitorId, input string, predicate ConnectedPartnerResourcesListFormatOperationPredicate) (result ConnectedPartnerResourcesListCompleteResult, err error) {
	items := make([]ConnectedPartnerResourcesListFormat, 0)

	resp, err := c.ConnectedPartnerResourcesList(ctx, id, input)
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

	result = ConnectedPartnerResourcesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
