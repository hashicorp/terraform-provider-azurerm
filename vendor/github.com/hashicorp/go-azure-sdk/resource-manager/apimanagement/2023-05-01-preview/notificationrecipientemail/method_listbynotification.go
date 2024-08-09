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

type ListByNotificationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecipientEmailContract
}

type ListByNotificationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RecipientEmailContract
}

type ListByNotificationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByNotificationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByNotification ...
func (c NotificationRecipientEmailClient) ListByNotification(ctx context.Context, id NotificationId) (result ListByNotificationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByNotificationCustomPager{},
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

// ListByNotificationComplete retrieves all the results into a single object
func (c NotificationRecipientEmailClient) ListByNotificationComplete(ctx context.Context, id NotificationId) (ListByNotificationCompleteResult, error) {
	return c.ListByNotificationCompleteMatchingPredicate(ctx, id, RecipientEmailContractOperationPredicate{})
}

// ListByNotificationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NotificationRecipientEmailClient) ListByNotificationCompleteMatchingPredicate(ctx context.Context, id NotificationId, predicate RecipientEmailContractOperationPredicate) (result ListByNotificationCompleteResult, err error) {
	items := make([]RecipientEmailContract, 0)

	resp, err := c.ListByNotification(ctx, id)
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

	result = ListByNotificationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
