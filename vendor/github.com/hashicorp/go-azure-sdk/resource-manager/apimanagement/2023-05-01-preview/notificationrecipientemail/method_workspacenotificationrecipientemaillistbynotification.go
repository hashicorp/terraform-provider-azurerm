package notificationrecipientemail

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceNotificationRecipientEmailListByNotificationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecipientEmailContract
}

type WorkspaceNotificationRecipientEmailListByNotificationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RecipientEmailContract
}

// WorkspaceNotificationRecipientEmailListByNotification ...
func (c NotificationRecipientEmailClient) WorkspaceNotificationRecipientEmailListByNotification(ctx context.Context, id NotificationNotificationId) (result WorkspaceNotificationRecipientEmailListByNotificationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/recipientEmails", id.ID()),
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
		Values *[]RecipientEmailContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceNotificationRecipientEmailListByNotificationComplete retrieves all the results into a single object
func (c NotificationRecipientEmailClient) WorkspaceNotificationRecipientEmailListByNotificationComplete(ctx context.Context, id NotificationNotificationId) (WorkspaceNotificationRecipientEmailListByNotificationCompleteResult, error) {
	return c.WorkspaceNotificationRecipientEmailListByNotificationCompleteMatchingPredicate(ctx, id, RecipientEmailContractOperationPredicate{})
}

// WorkspaceNotificationRecipientEmailListByNotificationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NotificationRecipientEmailClient) WorkspaceNotificationRecipientEmailListByNotificationCompleteMatchingPredicate(ctx context.Context, id NotificationNotificationId, predicate RecipientEmailContractOperationPredicate) (result WorkspaceNotificationRecipientEmailListByNotificationCompleteResult, err error) {
	items := make([]RecipientEmailContract, 0)

	resp, err := c.WorkspaceNotificationRecipientEmailListByNotification(ctx, id)
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

	result = WorkspaceNotificationRecipientEmailListByNotificationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
