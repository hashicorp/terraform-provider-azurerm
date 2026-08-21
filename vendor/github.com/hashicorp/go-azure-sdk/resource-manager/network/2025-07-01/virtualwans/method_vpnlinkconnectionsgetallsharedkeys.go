package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnLinkConnectionsGetAllSharedKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConnectionSharedKeyResult
}

type VpnLinkConnectionsGetAllSharedKeysCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConnectionSharedKeyResult
}

type VpnLinkConnectionsGetAllSharedKeysCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VpnLinkConnectionsGetAllSharedKeysCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VpnLinkConnectionsGetAllSharedKeys ...
func (c VirtualWANsClient) VpnLinkConnectionsGetAllSharedKeys(ctx context.Context, id VpnLinkConnectionId) (result VpnLinkConnectionsGetAllSharedKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VpnLinkConnectionsGetAllSharedKeysCustomPager{},
		Path:       fmt.Sprintf("%s/sharedKeys", id.ID()),
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
		Values *[]ConnectionSharedKeyResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnLinkConnectionsGetAllSharedKeysComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnLinkConnectionsGetAllSharedKeysComplete(ctx context.Context, id VpnLinkConnectionId) (VpnLinkConnectionsGetAllSharedKeysCompleteResult, error) {
	return c.VpnLinkConnectionsGetAllSharedKeysCompleteMatchingPredicate(ctx, id, ConnectionSharedKeyResultOperationPredicate{})
}

// VpnLinkConnectionsGetAllSharedKeysCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnLinkConnectionsGetAllSharedKeysCompleteMatchingPredicate(ctx context.Context, id VpnLinkConnectionId, predicate ConnectionSharedKeyResultOperationPredicate) (result VpnLinkConnectionsGetAllSharedKeysCompleteResult, err error) {
	items := make([]ConnectionSharedKeyResult, 0)

	resp, err := c.VpnLinkConnectionsGetAllSharedKeys(ctx, id)
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

	result = VpnLinkConnectionsGetAllSharedKeysCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
