package recordsets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAllByDnsZoneOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecordSet
}

type ListAllByDnsZoneCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RecordSet
}

type ListAllByDnsZoneOperationOptions struct {
	Recordsetnamesuffix *string
	Top                 *int64
}

func DefaultListAllByDnsZoneOperationOptions() ListAllByDnsZoneOperationOptions {
	return ListAllByDnsZoneOperationOptions{}
}

func (o ListAllByDnsZoneOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAllByDnsZoneOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAllByDnsZoneOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Recordsetnamesuffix != nil {
		out.Append("$recordsetnamesuffix", fmt.Sprintf("%v", *o.Recordsetnamesuffix))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListAllByDnsZoneCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAllByDnsZoneCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAllByDnsZone ...
func (c RecordSetsClient) ListAllByDnsZone(ctx context.Context, id DnsZoneId, options ListAllByDnsZoneOperationOptions) (result ListAllByDnsZoneOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAllByDnsZoneCustomPager{},
		Path:          fmt.Sprintf("%s/all", id.ID()),
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
		Values *[]RecordSet `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAllByDnsZoneComplete retrieves all the results into a single object
func (c RecordSetsClient) ListAllByDnsZoneComplete(ctx context.Context, id DnsZoneId, options ListAllByDnsZoneOperationOptions) (ListAllByDnsZoneCompleteResult, error) {
	return c.ListAllByDnsZoneCompleteMatchingPredicate(ctx, id, options, RecordSetOperationPredicate{})
}

// ListAllByDnsZoneCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RecordSetsClient) ListAllByDnsZoneCompleteMatchingPredicate(ctx context.Context, id DnsZoneId, options ListAllByDnsZoneOperationOptions, predicate RecordSetOperationPredicate) (result ListAllByDnsZoneCompleteResult, err error) {
	items := make([]RecordSet, 0)

	resp, err := c.ListAllByDnsZone(ctx, id, options)
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

	result = ListAllByDnsZoneCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
