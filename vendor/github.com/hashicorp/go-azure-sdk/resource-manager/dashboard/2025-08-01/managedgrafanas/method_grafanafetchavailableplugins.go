package managedgrafanas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GrafanaFetchAvailablePluginsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GrafanaAvailablePlugin
}

type GrafanaFetchAvailablePluginsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GrafanaAvailablePlugin
}

type GrafanaFetchAvailablePluginsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GrafanaFetchAvailablePluginsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GrafanaFetchAvailablePlugins ...
func (c ManagedGrafanasClient) GrafanaFetchAvailablePlugins(ctx context.Context, id GrafanaId) (result GrafanaFetchAvailablePluginsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &GrafanaFetchAvailablePluginsCustomPager{},
		Path:       fmt.Sprintf("%s/fetchAvailablePlugins", id.ID()),
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
		Values *[]GrafanaAvailablePlugin `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GrafanaFetchAvailablePluginsComplete retrieves all the results into a single object
func (c ManagedGrafanasClient) GrafanaFetchAvailablePluginsComplete(ctx context.Context, id GrafanaId) (GrafanaFetchAvailablePluginsCompleteResult, error) {
	return c.GrafanaFetchAvailablePluginsCompleteMatchingPredicate(ctx, id, GrafanaAvailablePluginOperationPredicate{})
}

// GrafanaFetchAvailablePluginsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedGrafanasClient) GrafanaFetchAvailablePluginsCompleteMatchingPredicate(ctx context.Context, id GrafanaId, predicate GrafanaAvailablePluginOperationPredicate) (result GrafanaFetchAvailablePluginsCompleteResult, err error) {
	items := make([]GrafanaAvailablePlugin, 0)

	resp, err := c.GrafanaFetchAvailablePlugins(ctx, id)
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

	result = GrafanaFetchAvailablePluginsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
