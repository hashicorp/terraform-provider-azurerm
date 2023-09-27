package webhooks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListEventsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Event
}

type ListEventsCompleteResult struct {
	Items []Event
}

// ListEvents ...
func (c WebHooksClient) ListEvents(ctx context.Context, id WebHookId) (result ListEventsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listEvents", id.ID()),
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
		Values *[]Event `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListEventsComplete retrieves all the results into a single object
func (c WebHooksClient) ListEventsComplete(ctx context.Context, id WebHookId) (ListEventsCompleteResult, error) {
	return c.ListEventsCompleteMatchingPredicate(ctx, id, EventOperationPredicate{})
}

// ListEventsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebHooksClient) ListEventsCompleteMatchingPredicate(ctx context.Context, id WebHookId, predicate EventOperationPredicate) (result ListEventsCompleteResult, err error) {
	items := make([]Event, 0)

	resp, err := c.ListEvents(ctx, id)
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

	result = ListEventsCompleteResult{
		Items: items,
	}
	return
}
