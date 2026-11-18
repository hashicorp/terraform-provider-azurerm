package connectionanalyzers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkWatchersConnectionAnalyzersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConnectionAnalyzer
}

type NetworkWatchersConnectionAnalyzersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConnectionAnalyzer
}

type NetworkWatchersConnectionAnalyzersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NetworkWatchersConnectionAnalyzersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NetworkWatchersConnectionAnalyzersList ...
func (c ConnectionAnalyzersClient) NetworkWatchersConnectionAnalyzersList(ctx context.Context, id NetworkWatcherId) (result NetworkWatchersConnectionAnalyzersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NetworkWatchersConnectionAnalyzersListCustomPager{},
		Path:       fmt.Sprintf("%s/connectionAnalyzers", id.ID()),
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
		Values *[]ConnectionAnalyzer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NetworkWatchersConnectionAnalyzersListComplete retrieves all the results into a single object
func (c ConnectionAnalyzersClient) NetworkWatchersConnectionAnalyzersListComplete(ctx context.Context, id NetworkWatcherId) (NetworkWatchersConnectionAnalyzersListCompleteResult, error) {
	return c.NetworkWatchersConnectionAnalyzersListCompleteMatchingPredicate(ctx, id, ConnectionAnalyzerOperationPredicate{})
}

// NetworkWatchersConnectionAnalyzersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConnectionAnalyzersClient) NetworkWatchersConnectionAnalyzersListCompleteMatchingPredicate(ctx context.Context, id NetworkWatcherId, predicate ConnectionAnalyzerOperationPredicate) (result NetworkWatchersConnectionAnalyzersListCompleteResult, err error) {
	items := make([]ConnectionAnalyzer, 0)

	resp, err := c.NetworkWatchersConnectionAnalyzersList(ctx, id)
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

	result = NetworkWatchersConnectionAnalyzersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
