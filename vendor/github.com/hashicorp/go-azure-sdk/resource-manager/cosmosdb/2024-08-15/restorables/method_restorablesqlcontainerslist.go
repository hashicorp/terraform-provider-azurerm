package restorables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableSqlContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RestorableSqlContainersListResult
}

type RestorableSqlContainersListOperationOptions struct {
	EndTime                  *string
	RestorableSqlDatabaseRid *string
	StartTime                *string
}

func DefaultRestorableSqlContainersListOperationOptions() RestorableSqlContainersListOperationOptions {
	return RestorableSqlContainersListOperationOptions{}
}

func (o RestorableSqlContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableSqlContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableSqlContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.RestorableSqlDatabaseRid != nil {
		out.Append("restorableSqlDatabaseRid", fmt.Sprintf("%v", *o.RestorableSqlDatabaseRid))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

// RestorableSqlContainersList ...
func (c RestorablesClient) RestorableSqlContainersList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlContainersListOperationOptions) (result RestorableSqlContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restorableSqlContainers", id.ID()),
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

	var model RestorableSqlContainersListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
