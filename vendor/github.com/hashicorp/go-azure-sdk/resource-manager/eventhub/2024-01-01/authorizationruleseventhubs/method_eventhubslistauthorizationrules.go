package authorizationruleseventhubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubsListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AuthorizationRule
}

type EventHubsListAuthorizationRulesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AuthorizationRule
}

type EventHubsListAuthorizationRulesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *EventHubsListAuthorizationRulesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// EventHubsListAuthorizationRules ...
func (c AuthorizationRulesEventHubsClient) EventHubsListAuthorizationRules(ctx context.Context, id EventhubId) (result EventHubsListAuthorizationRulesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &EventHubsListAuthorizationRulesCustomPager{},
		Path:       fmt.Sprintf("%s/authorizationRules", id.ID()),
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
		Values *[]AuthorizationRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// EventHubsListAuthorizationRulesComplete retrieves all the results into a single object
func (c AuthorizationRulesEventHubsClient) EventHubsListAuthorizationRulesComplete(ctx context.Context, id EventhubId) (EventHubsListAuthorizationRulesCompleteResult, error) {
	return c.EventHubsListAuthorizationRulesCompleteMatchingPredicate(ctx, id, AuthorizationRuleOperationPredicate{})
}

// EventHubsListAuthorizationRulesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AuthorizationRulesEventHubsClient) EventHubsListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id EventhubId, predicate AuthorizationRuleOperationPredicate) (result EventHubsListAuthorizationRulesCompleteResult, err error) {
	items := make([]AuthorizationRule, 0)

	resp, err := c.EventHubsListAuthorizationRules(ctx, id)
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

	result = EventHubsListAuthorizationRulesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
