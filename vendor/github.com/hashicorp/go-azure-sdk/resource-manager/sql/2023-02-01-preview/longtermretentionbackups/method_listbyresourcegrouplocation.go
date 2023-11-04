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

type ListByResourceGroupLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LongTermRetentionBackup
}

type ListByResourceGroupLocationCompleteResult struct {
	Items []LongTermRetentionBackup
}

type ListByResourceGroupLocationOperationOptions struct {
	DatabaseState         *DatabaseState
	OnlyLatestPerDatabase *bool
}

func DefaultListByResourceGroupLocationOperationOptions() ListByResourceGroupLocationOperationOptions {
	return ListByResourceGroupLocationOperationOptions{}
}

func (o ListByResourceGroupLocationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByResourceGroupLocationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByResourceGroupLocationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DatabaseState != nil {
		out.Append("databaseState", fmt.Sprintf("%v", *o.DatabaseState))
	}
	if o.OnlyLatestPerDatabase != nil {
		out.Append("onlyLatestPerDatabase", fmt.Sprintf("%v", *o.OnlyLatestPerDatabase))
	}
	return &out
}

// ListByResourceGroupLocation ...
func (c LongTermRetentionBackupsClient) ListByResourceGroupLocation(ctx context.Context, id ProviderLocationId, options ListByResourceGroupLocationOperationOptions) (result ListByResourceGroupLocationOperationResponse, err error) {
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

// ListByResourceGroupLocationComplete retrieves all the results into a single object
func (c LongTermRetentionBackupsClient) ListByResourceGroupLocationComplete(ctx context.Context, id ProviderLocationId, options ListByResourceGroupLocationOperationOptions) (ListByResourceGroupLocationCompleteResult, error) {
	return c.ListByResourceGroupLocationCompleteMatchingPredicate(ctx, id, options, LongTermRetentionBackupOperationPredicate{})
}

// ListByResourceGroupLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LongTermRetentionBackupsClient) ListByResourceGroupLocationCompleteMatchingPredicate(ctx context.Context, id ProviderLocationId, options ListByResourceGroupLocationOperationOptions, predicate LongTermRetentionBackupOperationPredicate) (result ListByResourceGroupLocationCompleteResult, err error) {
	items := make([]LongTermRetentionBackup, 0)

	resp, err := c.ListByResourceGroupLocation(ctx, id, options)
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

	result = ListByResourceGroupLocationCompleteResult{
		Items: items,
	}
	return
}
