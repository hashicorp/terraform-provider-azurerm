package longtermretentionbackups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LongTermRetentionBackup
}

type ListByServerCompleteResult struct {
	Items []LongTermRetentionBackup
}

type ListByServerOperationOptions struct {
	DatabaseState         *DatabaseState
	OnlyLatestPerDatabase *bool
}

func DefaultListByServerOperationOptions() ListByServerOperationOptions {
	return ListByServerOperationOptions{}
}

func (o ListByServerOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByServerOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByServerOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DatabaseState != nil {
		out.Append("databaseState", fmt.Sprintf("%v", *o.DatabaseState))
	}
	if o.OnlyLatestPerDatabase != nil {
		out.Append("onlyLatestPerDatabase", fmt.Sprintf("%v", *o.OnlyLatestPerDatabase))
	}
	return &out
}

// ListByServer ...
func (c LongTermRetentionBackupsClient) ListByServer(ctx context.Context, id LongTermRetentionServerId, options ListByServerOperationOptions) (result ListByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/longTermRetentionBackups", id.ID()),
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
		Values *[]LongTermRetentionBackup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByServerComplete retrieves all the results into a single object
func (c LongTermRetentionBackupsClient) ListByServerComplete(ctx context.Context, id LongTermRetentionServerId, options ListByServerOperationOptions) (ListByServerCompleteResult, error) {
	return c.ListByServerCompleteMatchingPredicate(ctx, id, options, LongTermRetentionBackupOperationPredicate{})
}

// ListByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LongTermRetentionBackupsClient) ListByServerCompleteMatchingPredicate(ctx context.Context, id LongTermRetentionServerId, options ListByServerOperationOptions, predicate LongTermRetentionBackupOperationPredicate) (result ListByServerCompleteResult, err error) {
	items := make([]LongTermRetentionBackup, 0)

	resp, err := c.ListByServer(ctx, id, options)
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

	result = ListByServerCompleteResult{
		Items: items,
	}
	return
}
