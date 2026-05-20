package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CollectionPartitionListUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PartitionUsagesResult
}

type CollectionPartitionListUsagesOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionListUsagesOperationOptions() CollectionPartitionListUsagesOperationOptions {
	return CollectionPartitionListUsagesOperationOptions{}
}

func (o CollectionPartitionListUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionPartitionListUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CollectionPartitionListUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// CollectionPartitionListUsages ...
func (c CosmosDBClient) CollectionPartitionListUsages(ctx context.Context, id CollectionId, options CollectionPartitionListUsagesOperationOptions) (result CollectionPartitionListUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/partitions/usages", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model PartitionUsagesResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
