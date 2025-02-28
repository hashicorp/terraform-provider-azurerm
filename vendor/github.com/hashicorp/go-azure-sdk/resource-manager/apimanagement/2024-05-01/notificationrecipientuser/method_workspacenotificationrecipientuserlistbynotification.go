package notificationrecipientuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceNotificationRecipientUserListByNotificationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecipientUserContract
}

type WorkspaceNotificationRecipientUserListByNotificationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RecipientUserContract
}

type WorkspaceNotificationRecipientUserListByNotificationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceNotificationRecipientUserListByNotificationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceNotificationRecipientUserListByNotification ...
func (c NotificationRecipientUserClient) WorkspaceNotificationRecipientUserListByNotification(ctx context.Context, id NotificationNotificationId) (result WorkspaceNotificationRecipientUserListByNotificationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &WorkspaceNotificationRecipientUserListByNotificationCustomPager{},
		Path:       fmt.Sprintf("%s/recipientUsers", id.ID()),
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
		Values *[]RecipientUserContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceNotificationRecipientUserListByNotificationComplete retrieves all the results into a single object
func (c NotificationRecipientUserClient) WorkspaceNotificationRecipientUserListByNotificationComplete(ctx context.Context, id NotificationNotificationId) (WorkspaceNotificationRecipientUserListByNotificationCompleteResult, error) {
	return c.WorkspaceNotificationRecipientUserListByNotificationCompleteMatchingPredicate(ctx, id, RecipientUserContractOperationPredicate{})
}

// WorkspaceNotificationRecipientUserListByNotificationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NotificationRecipientUserClient) WorkspaceNotificationRecipientUserListByNotificationCompleteMatchingPredicate(ctx context.Context, id NotificationNotificationId, predicate RecipientUserContractOperationPredicate) (result WorkspaceNotificationRecipientUserListByNotificationCompleteResult, err error) {
	items := make([]RecipientUserContract, 0)

	resp, err := c.WorkspaceNotificationRecipientUserListByNotification(ctx, id)
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

	result = WorkspaceNotificationRecipientUserListByNotificationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
