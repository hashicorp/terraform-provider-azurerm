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

type ListByDnsZoneOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecordSet
}

type ListByDnsZoneCompleteResult struct {
	Items []RecordSet
}

type ListByDnsZoneOperationOptions struct {
	Recordsetnamesuffix *string
	Top                 *int64
}

func DefaultListByDnsZoneOperationOptions() ListByDnsZoneOperationOptions {
	return ListByDnsZoneOperationOptions{}
}

func (o ListByDnsZoneOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByDnsZoneOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByDnsZoneOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Recordsetnamesuffix != nil {
		out.Append("$recordsetnamesuffix", fmt.Sprintf("%v", *o.Recordsetnamesuffix))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByDnsZone ...
func (c RecordSetsClient) ListByDnsZone(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions) (result ListByDnsZoneOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/recordsets", id.ID()),
		OptionsObject: options,
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

// ListByDnsZoneComplete retrieves all the results into a single object
func (c RecordSetsClient) ListByDnsZoneComplete(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions) (ListByDnsZoneCompleteResult, error) {
	return c.ListByDnsZoneCompleteMatchingPredicate(ctx, id, options, RecordSetOperationPredicate{})
}

// ListByDnsZoneCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RecordSetsClient) ListByDnsZoneCompleteMatchingPredicate(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions, predicate RecordSetOperationPredicate) (result ListByDnsZoneCompleteResult, err error) {
	items := make([]RecordSet, 0)

	resp, err := c.ListByDnsZone(ctx, id, options)
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

	result = ListByDnsZoneCompleteResult{
		Items: items,
	}
	return
}
