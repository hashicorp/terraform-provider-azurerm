package privatedns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSetsListByTypeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecordSet
}

type RecordSetsListByTypeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RecordSet
}

type RecordSetsListByTypeOperationOptions struct {
	Recordsetnamesuffix *string
	Top                 *int64
}

func DefaultRecordSetsListByTypeOperationOptions() RecordSetsListByTypeOperationOptions {
	return RecordSetsListByTypeOperationOptions{}
}

func (o RecordSetsListByTypeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RecordSetsListByTypeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RecordSetsListByTypeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Recordsetnamesuffix != nil {
		out.Append("$recordsetnamesuffix", fmt.Sprintf("%v", *o.Recordsetnamesuffix))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type RecordSetsListByTypeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RecordSetsListByTypeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RecordSetsListByType ...
func (c PrivateDNSClient) RecordSetsListByType(ctx context.Context, id PrivateZoneId, options RecordSetsListByTypeOperationOptions) (result RecordSetsListByTypeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RecordSetsListByTypeCustomPager{},
		Path:          id.ID(),
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

// RecordSetsListByTypeComplete retrieves all the results into a single object
func (c PrivateDNSClient) RecordSetsListByTypeComplete(ctx context.Context, id PrivateZoneId, options RecordSetsListByTypeOperationOptions) (RecordSetsListByTypeCompleteResult, error) {
	return c.RecordSetsListByTypeCompleteMatchingPredicate(ctx, id, options, RecordSetOperationPredicate{})
}

// RecordSetsListByTypeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateDNSClient) RecordSetsListByTypeCompleteMatchingPredicate(ctx context.Context, id PrivateZoneId, options RecordSetsListByTypeOperationOptions, predicate RecordSetOperationPredicate) (result RecordSetsListByTypeCompleteResult, err error) {
	items := make([]RecordSet, 0)

	resp, err := c.RecordSetsListByType(ctx, id, options)
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

	result = RecordSetsListByTypeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
