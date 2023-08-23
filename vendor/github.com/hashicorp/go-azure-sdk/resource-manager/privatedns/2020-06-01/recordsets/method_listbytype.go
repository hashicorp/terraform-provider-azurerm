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

type ListByTypeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecordSet
}

type ListByTypeCompleteResult struct {
	Items []RecordSet
}

type ListByTypeOperationOptions struct {
	Recordsetnamesuffix *string
	Top                 *int64
}

func DefaultListByTypeOperationOptions() ListByTypeOperationOptions {
	return ListByTypeOperationOptions{}
}

func (o ListByTypeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByTypeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByTypeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Recordsetnamesuffix != nil {
		out.Append("$recordsetnamesuffix", fmt.Sprintf("%v", *o.Recordsetnamesuffix))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByType ...
func (c RecordSetsClient) ListByType(ctx context.Context, id PrivateZoneId, options ListByTypeOperationOptions) (result ListByTypeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          id.ID(),
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

// ListByTypeComplete retrieves all the results into a single object
func (c RecordSetsClient) ListByTypeComplete(ctx context.Context, id PrivateZoneId, options ListByTypeOperationOptions) (ListByTypeCompleteResult, error) {
	return c.ListByTypeCompleteMatchingPredicate(ctx, id, options, RecordSetOperationPredicate{})
}

// ListByTypeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RecordSetsClient) ListByTypeCompleteMatchingPredicate(ctx context.Context, id PrivateZoneId, options ListByTypeOperationOptions, predicate RecordSetOperationPredicate) (result ListByTypeCompleteResult, err error) {
	items := make([]RecordSet, 0)

	resp, err := c.ListByType(ctx, id, options)
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

	result = ListByTypeCompleteResult{
		Items: items,
	}
	return
}
