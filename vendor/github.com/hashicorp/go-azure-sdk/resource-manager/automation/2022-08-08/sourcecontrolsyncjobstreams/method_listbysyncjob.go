package sourcecontrolsyncjobstreams

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySyncJobOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SourceControlSyncJobStream
}

type ListBySyncJobCompleteResult struct {
	Items []SourceControlSyncJobStream
}

type ListBySyncJobOperationOptions struct {
	Filter *string
}

func DefaultListBySyncJobOperationOptions() ListBySyncJobOperationOptions {
	return ListBySyncJobOperationOptions{}
}

func (o ListBySyncJobOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListBySyncJobOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListBySyncJobOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ListBySyncJob ...
func (c SourceControlSyncJobStreamsClient) ListBySyncJob(ctx context.Context, id SourceControlSyncJobId, options ListBySyncJobOperationOptions) (result ListBySyncJobOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/streams", id.ID()),
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
		Values *[]SourceControlSyncJobStream `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySyncJobComplete retrieves all the results into a single object
func (c SourceControlSyncJobStreamsClient) ListBySyncJobComplete(ctx context.Context, id SourceControlSyncJobId, options ListBySyncJobOperationOptions) (ListBySyncJobCompleteResult, error) {
	return c.ListBySyncJobCompleteMatchingPredicate(ctx, id, options, SourceControlSyncJobStreamOperationPredicate{})
}

// ListBySyncJobCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SourceControlSyncJobStreamsClient) ListBySyncJobCompleteMatchingPredicate(ctx context.Context, id SourceControlSyncJobId, options ListBySyncJobOperationOptions, predicate SourceControlSyncJobStreamOperationPredicate) (result ListBySyncJobCompleteResult, err error) {
	items := make([]SourceControlSyncJobStream, 0)

	resp, err := c.ListBySyncJob(ctx, id, options)
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

	result = ListBySyncJobCompleteResult{
		Items: items,
	}
	return
}
