package querykeys

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySearchServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]QueryKey
}

type ListBySearchServiceCompleteResult struct {
	Items []QueryKey
}

type ListBySearchServiceOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultListBySearchServiceOperationOptions() ListBySearchServiceOperationOptions {
	return ListBySearchServiceOperationOptions{}
}

func (o ListBySearchServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o ListBySearchServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListBySearchServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// ListBySearchService ...
func (c QueryKeysClient) ListBySearchService(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions) (result ListBySearchServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/listQueryKeys", id.ID()),
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
		Values *[]QueryKey `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySearchServiceComplete retrieves all the results into a single object
func (c QueryKeysClient) ListBySearchServiceComplete(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions) (ListBySearchServiceCompleteResult, error) {
	return c.ListBySearchServiceCompleteMatchingPredicate(ctx, id, options, QueryKeyOperationPredicate{})
}

// ListBySearchServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c QueryKeysClient) ListBySearchServiceCompleteMatchingPredicate(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions, predicate QueryKeyOperationPredicate) (result ListBySearchServiceCompleteResult, err error) {
	items := make([]QueryKey, 0)

	resp, err := c.ListBySearchService(ctx, id, options)
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

	result = ListBySearchServiceCompleteResult{
		Items: items,
	}
	return
}
