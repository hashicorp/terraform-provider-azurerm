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

type ListByResourceGroupServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LongTermRetentionBackup
}

type ListByResourceGroupServerCompleteResult struct {
	Items []LongTermRetentionBackup
}

type ListByResourceGroupServerOperationOptions struct {
	DatabaseState         *DatabaseState
	OnlyLatestPerDatabase *bool
}

func DefaultListByResourceGroupServerOperationOptions() ListByResourceGroupServerOperationOptions {
	return ListByResourceGroupServerOperationOptions{}
}

func (o ListByResourceGroupServerOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByResourceGroupServerOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByResourceGroupServerOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DatabaseState != nil {
		out.Append("databaseState", fmt.Sprintf("%v", *o.DatabaseState))
	}
	if o.OnlyLatestPerDatabase != nil {
		out.Append("onlyLatestPerDatabase", fmt.Sprintf("%v", *o.OnlyLatestPerDatabase))
	}
	return &out
}

// ListByResourceGroupServer ...
func (c LongTermRetentionBackupsClient) ListByResourceGroupServer(ctx context.Context, id LocationLongTermRetentionServerId, options ListByResourceGroupServerOperationOptions) (result ListByResourceGroupServerOperationResponse, err error) {
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

// ListByResourceGroupServerComplete retrieves all the results into a single object
func (c LongTermRetentionBackupsClient) ListByResourceGroupServerComplete(ctx context.Context, id LocationLongTermRetentionServerId, options ListByResourceGroupServerOperationOptions) (ListByResourceGroupServerCompleteResult, error) {
	return c.ListByResourceGroupServerCompleteMatchingPredicate(ctx, id, options, LongTermRetentionBackupOperationPredicate{})
}

// ListByResourceGroupServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LongTermRetentionBackupsClient) ListByResourceGroupServerCompleteMatchingPredicate(ctx context.Context, id LocationLongTermRetentionServerId, options ListByResourceGroupServerOperationOptions, predicate LongTermRetentionBackupOperationPredicate) (result ListByResourceGroupServerCompleteResult, err error) {
	items := make([]LongTermRetentionBackup, 0)

	resp, err := c.ListByResourceGroupServer(ctx, id, options)
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

	result = ListByResourceGroupServerCompleteResult{
		Items: items,
	}
	return
}
