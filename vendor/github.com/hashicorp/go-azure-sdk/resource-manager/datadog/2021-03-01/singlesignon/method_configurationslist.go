package singlesignon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatadogSingleSignOnResource
}

type ConfigurationsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatadogSingleSignOnResource
}

type ConfigurationsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ConfigurationsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ConfigurationsList ...
func (c SingleSignOnClient) ConfigurationsList(ctx context.Context, id MonitorId) (result ConfigurationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ConfigurationsListCustomPager{},
		Path:       fmt.Sprintf("%s/singleSignOnConfigurations", id.ID()),
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
		Values *[]DatadogSingleSignOnResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ConfigurationsListComplete retrieves all the results into a single object
func (c SingleSignOnClient) ConfigurationsListComplete(ctx context.Context, id MonitorId) (ConfigurationsListCompleteResult, error) {
	return c.ConfigurationsListCompleteMatchingPredicate(ctx, id, DatadogSingleSignOnResourceOperationPredicate{})
}

// ConfigurationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SingleSignOnClient) ConfigurationsListCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate DatadogSingleSignOnResourceOperationPredicate) (result ConfigurationsListCompleteResult, err error) {
	items := make([]DatadogSingleSignOnResource, 0)

	resp, err := c.ConfigurationsList(ctx, id)
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

	result = ConfigurationsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
