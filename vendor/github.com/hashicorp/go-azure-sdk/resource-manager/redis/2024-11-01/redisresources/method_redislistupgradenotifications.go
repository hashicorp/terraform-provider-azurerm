package redisresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisListUpgradeNotificationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UpgradeNotification
}

type RedisListUpgradeNotificationsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []UpgradeNotification
}

type RedisListUpgradeNotificationsOperationOptions struct {
	History *float64
}

func DefaultRedisListUpgradeNotificationsOperationOptions() RedisListUpgradeNotificationsOperationOptions {
	return RedisListUpgradeNotificationsOperationOptions{}
}

func (o RedisListUpgradeNotificationsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RedisListUpgradeNotificationsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RedisListUpgradeNotificationsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.History != nil {
		out.Append("history", fmt.Sprintf("%v", *o.History))
	}
	return &out
}

type RedisListUpgradeNotificationsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RedisListUpgradeNotificationsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RedisListUpgradeNotifications ...
func (c RedisResourcesClient) RedisListUpgradeNotifications(ctx context.Context, id RediId, options RedisListUpgradeNotificationsOperationOptions) (result RedisListUpgradeNotificationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RedisListUpgradeNotificationsCustomPager{},
		Path:          fmt.Sprintf("%s/listUpgradeNotifications", id.ID()),
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
		Values *[]UpgradeNotification `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RedisListUpgradeNotificationsComplete retrieves all the results into a single object
func (c RedisResourcesClient) RedisListUpgradeNotificationsComplete(ctx context.Context, id RediId, options RedisListUpgradeNotificationsOperationOptions) (RedisListUpgradeNotificationsCompleteResult, error) {
	return c.RedisListUpgradeNotificationsCompleteMatchingPredicate(ctx, id, options, UpgradeNotificationOperationPredicate{})
}

// RedisListUpgradeNotificationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisResourcesClient) RedisListUpgradeNotificationsCompleteMatchingPredicate(ctx context.Context, id RediId, options RedisListUpgradeNotificationsOperationOptions, predicate UpgradeNotificationOperationPredicate) (result RedisListUpgradeNotificationsCompleteResult, err error) {
	items := make([]UpgradeNotification, 0)

	resp, err := c.RedisListUpgradeNotifications(ctx, id, options)
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

	result = RedisListUpgradeNotificationsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
