package apikey

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsListApiKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatadogApiKey
}

type MonitorsListApiKeysCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatadogApiKey
}

type MonitorsListApiKeysCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListApiKeysCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListApiKeys ...
func (c ApiKeyClient) MonitorsListApiKeys(ctx context.Context, id MonitorId) (result MonitorsListApiKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &MonitorsListApiKeysCustomPager{},
		Path:       fmt.Sprintf("%s/listApiKeys", id.ID()),
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
		Values *[]DatadogApiKey `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MonitorsListApiKeysComplete retrieves all the results into a single object
func (c ApiKeyClient) MonitorsListApiKeysComplete(ctx context.Context, id MonitorId) (MonitorsListApiKeysCompleteResult, error) {
	return c.MonitorsListApiKeysCompleteMatchingPredicate(ctx, id, DatadogApiKeyOperationPredicate{})
}

// MonitorsListApiKeysCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiKeyClient) MonitorsListApiKeysCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate DatadogApiKeyOperationPredicate) (result MonitorsListApiKeysCompleteResult, err error) {
	items := make([]DatadogApiKey, 0)

	resp, err := c.MonitorsListApiKeys(ctx, id)
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

	result = MonitorsListApiKeysCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
