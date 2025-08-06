package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainersListLogsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Logs
}

type ContainersListLogsOperationOptions struct {
	Tail       *int64
	Timestamps *bool
}

func DefaultContainersListLogsOperationOptions() ContainersListLogsOperationOptions {
	return ContainersListLogsOperationOptions{}
}

func (o ContainersListLogsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ContainersListLogsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ContainersListLogsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Tail != nil {
		out.Append("tail", fmt.Sprintf("%v", *o.Tail))
	}
	if o.Timestamps != nil {
		out.Append("timestamps", fmt.Sprintf("%v", *o.Timestamps))
	}
	return &out
}

// ContainersListLogs ...
func (c ContainerInstanceClient) ContainersListLogs(ctx context.Context, id ContainerId, options ContainersListLogsOperationOptions) (result ContainersListLogsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/logs", id.ID()),
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

	var model Logs
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
