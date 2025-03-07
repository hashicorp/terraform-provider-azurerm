package notification

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceNotificationListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NotificationContract
}

type WorkspaceNotificationListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NotificationContract
}

type WorkspaceNotificationListByServiceOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultWorkspaceNotificationListByServiceOperationOptions() WorkspaceNotificationListByServiceOperationOptions {
	return WorkspaceNotificationListByServiceOperationOptions{}
}

func (o WorkspaceNotificationListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceNotificationListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceNotificationListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceNotificationListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceNotificationListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceNotificationListByService ...
func (c NotificationClient) WorkspaceNotificationListByService(ctx context.Context, id WorkspaceId, options WorkspaceNotificationListByServiceOperationOptions) (result WorkspaceNotificationListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceNotificationListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/notifications", id.ID()),
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
		Values *[]NotificationContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceNotificationListByServiceComplete retrieves all the results into a single object
func (c NotificationClient) WorkspaceNotificationListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspaceNotificationListByServiceOperationOptions) (WorkspaceNotificationListByServiceCompleteResult, error) {
	return c.WorkspaceNotificationListByServiceCompleteMatchingPredicate(ctx, id, options, NotificationContractOperationPredicate{})
}

// WorkspaceNotificationListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NotificationClient) WorkspaceNotificationListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceNotificationListByServiceOperationOptions, predicate NotificationContractOperationPredicate) (result WorkspaceNotificationListByServiceCompleteResult, err error) {
	items := make([]NotificationContract, 0)

	resp, err := c.WorkspaceNotificationListByService(ctx, id, options)
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

	result = WorkspaceNotificationListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
