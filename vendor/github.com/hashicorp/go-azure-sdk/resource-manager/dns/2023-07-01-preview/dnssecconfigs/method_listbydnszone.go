package dnssecconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDnsZoneOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DnssecConfig
}

type ListByDnsZoneCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DnssecConfig
}

type ListByDnsZoneCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDnsZoneCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDnsZone ...
func (c DnssecConfigsClient) ListByDnsZone(ctx context.Context, id DnsZoneId) (result ListByDnsZoneOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByDnsZoneCustomPager{},
		Path:       fmt.Sprintf("%s/dnssecConfigs", id.ID()),
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
		Values *[]DnssecConfig `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDnsZoneComplete retrieves all the results into a single object
func (c DnssecConfigsClient) ListByDnsZoneComplete(ctx context.Context, id DnsZoneId) (ListByDnsZoneCompleteResult, error) {
	return c.ListByDnsZoneCompleteMatchingPredicate(ctx, id, DnssecConfigOperationPredicate{})
}

// ListByDnsZoneCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DnssecConfigsClient) ListByDnsZoneCompleteMatchingPredicate(ctx context.Context, id DnsZoneId, predicate DnssecConfigOperationPredicate) (result ListByDnsZoneCompleteResult, err error) {
	items := make([]DnssecConfig, 0)

	resp, err := c.ListByDnsZone(ctx, id)
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

	result = ListByDnsZoneCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
