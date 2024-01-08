package redis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListUpgradeNotificationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UpgradeNotification
}

type ListUpgradeNotificationsCompleteResult struct {
	Items []UpgradeNotification
}

type ListUpgradeNotificationsOperationOptions struct {
	History *float64
}

func DefaultListUpgradeNotificationsOperationOptions() ListUpgradeNotificationsOperationOptions {
	return ListUpgradeNotificationsOperationOptions{}
}

func (o ListUpgradeNotificationsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListUpgradeNotificationsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListUpgradeNotificationsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.History != nil {
		out.Append("history", fmt.Sprintf("%v", *o.History))
	}
	return &out
}

// ListUpgradeNotifications ...
func (c RedisClient) ListUpgradeNotifications(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions) (result ListUpgradeNotificationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/listUpgradeNotifications", id.ID()),
		OptionsObject: options,
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

// ListUpgradeNotificationsComplete retrieves all the results into a single object
func (c RedisClient) ListUpgradeNotificationsComplete(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions) (ListUpgradeNotificationsCompleteResult, error) {
	return c.ListUpgradeNotificationsCompleteMatchingPredicate(ctx, id, options, UpgradeNotificationOperationPredicate{})
}

// ListUpgradeNotificationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisClient) ListUpgradeNotificationsCompleteMatchingPredicate(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions, predicate UpgradeNotificationOperationPredicate) (result ListUpgradeNotificationsCompleteResult, err error) {
	items := make([]UpgradeNotification, 0)

	resp, err := c.ListUpgradeNotifications(ctx, id, options)
	if err != nil {
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

	result = ListUpgradeNotificationsCompleteResult{
		Items: items,
	}
	return
}
