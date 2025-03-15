package notification

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceNotificationCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NotificationContract
}

type WorkspaceNotificationCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceNotificationCreateOrUpdateOperationOptions() WorkspaceNotificationCreateOrUpdateOperationOptions {
	return WorkspaceNotificationCreateOrUpdateOperationOptions{}
}

func (o WorkspaceNotificationCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceNotificationCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceNotificationCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceNotificationCreateOrUpdate ...
func (c NotificationClient) WorkspaceNotificationCreateOrUpdate(ctx context.Context, id NotificationNotificationId, options WorkspaceNotificationCreateOrUpdateOperationOptions) (result WorkspaceNotificationCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          id.ID(),
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

	var model NotificationContract
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
