package bastionhosts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisconnectActiveSessionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BastionSessionState
}

type DisconnectActiveSessionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BastionSessionState
}

type DisconnectActiveSessionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DisconnectActiveSessionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DisconnectActiveSessions ...
func (c BastionHostsClient) DisconnectActiveSessions(ctx context.Context, id BastionHostId, input SessionIds) (result DisconnectActiveSessionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &DisconnectActiveSessionsCustomPager{},
		Path:       fmt.Sprintf("%s/disconnectActiveSessions", id.ID()),
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
		Values *[]BastionSessionState `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DisconnectActiveSessionsComplete retrieves all the results into a single object
func (c BastionHostsClient) DisconnectActiveSessionsComplete(ctx context.Context, id BastionHostId, input SessionIds) (DisconnectActiveSessionsCompleteResult, error) {
	return c.DisconnectActiveSessionsCompleteMatchingPredicate(ctx, id, input, BastionSessionStateOperationPredicate{})
}

// DisconnectActiveSessionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BastionHostsClient) DisconnectActiveSessionsCompleteMatchingPredicate(ctx context.Context, id BastionHostId, input SessionIds, predicate BastionSessionStateOperationPredicate) (result DisconnectActiveSessionsCompleteResult, err error) {
	items := make([]BastionSessionState, 0)

	resp, err := c.DisconnectActiveSessions(ctx, id, input)
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

	result = DisconnectActiveSessionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
