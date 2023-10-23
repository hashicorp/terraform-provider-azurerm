package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListLinkableEnvironmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LinkableEnvironmentResponse
}

type ListLinkableEnvironmentsCompleteResult struct {
	Items []LinkableEnvironmentResponse
}

// ListLinkableEnvironments ...
func (c MonitorsClient) ListLinkableEnvironments(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest) (result ListLinkableEnvironmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listLinkableEnvironments", id.ID()),
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
		Values *[]LinkableEnvironmentResponse `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListLinkableEnvironmentsComplete retrieves all the results into a single object
func (c MonitorsClient) ListLinkableEnvironmentsComplete(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest) (ListLinkableEnvironmentsCompleteResult, error) {
	return c.ListLinkableEnvironmentsCompleteMatchingPredicate(ctx, id, input, LinkableEnvironmentResponseOperationPredicate{})
}

// ListLinkableEnvironmentsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsClient) ListLinkableEnvironmentsCompleteMatchingPredicate(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest, predicate LinkableEnvironmentResponseOperationPredicate) (result ListLinkableEnvironmentsCompleteResult, err error) {
	items := make([]LinkableEnvironmentResponse, 0)

	resp, err := c.ListLinkableEnvironments(ctx, id, input)
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

	result = ListLinkableEnvironmentsCompleteResult{
		Items: items,
	}
	return
}
