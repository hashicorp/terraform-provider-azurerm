package hubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationHubsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NotificationHubResource
}

type NotificationHubsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NotificationHubResource
}

type NotificationHubsListOperationOptions struct {
	Top *int64
}

func DefaultNotificationHubsListOperationOptions() NotificationHubsListOperationOptions {
	return NotificationHubsListOperationOptions{}
}

func (o NotificationHubsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o NotificationHubsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o NotificationHubsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type NotificationHubsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NotificationHubsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NotificationHubsList ...
func (c HubsClient) NotificationHubsList(ctx context.Context, id NamespaceId, options NotificationHubsListOperationOptions) (result NotificationHubsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &NotificationHubsListCustomPager{},
		Path:          fmt.Sprintf("%s/notificationHubs", id.ID()),
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
		Values *[]NotificationHubResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NotificationHubsListComplete retrieves all the results into a single object
func (c HubsClient) NotificationHubsListComplete(ctx context.Context, id NamespaceId, options NotificationHubsListOperationOptions) (NotificationHubsListCompleteResult, error) {
	return c.NotificationHubsListCompleteMatchingPredicate(ctx, id, options, NotificationHubResourceOperationPredicate{})
}

// NotificationHubsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HubsClient) NotificationHubsListCompleteMatchingPredicate(ctx context.Context, id NamespaceId, options NotificationHubsListOperationOptions, predicate NotificationHubResourceOperationPredicate) (result NotificationHubsListCompleteResult, err error) {
	items := make([]NotificationHubResource, 0)

	resp, err := c.NotificationHubsList(ctx, id, options)
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

	result = NotificationHubsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
