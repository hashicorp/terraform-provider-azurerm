package backupinstanceresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupInstancesDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type BackupInstancesDeleteOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultBackupInstancesDeleteOperationOptions() BackupInstancesDeleteOperationOptions {
	return BackupInstancesDeleteOperationOptions{}
}

func (o BackupInstancesDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsAuthorizationAuxiliary != nil {
		out.Append("x-ms-authorization-auxiliary", fmt.Sprintf("%v", *o.XMsAuthorizationAuxiliary))
	}
	return &out
}

func (o BackupInstancesDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o BackupInstancesDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// BackupInstancesDelete ...
func (c BackupInstanceResourcesClient) BackupInstancesDelete(ctx context.Context, id BackupInstanceId, options BackupInstancesDeleteOperationOptions) (result BackupInstancesDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// BackupInstancesDeleteThenPoll performs BackupInstancesDelete then polls until it's completed
func (c BackupInstanceResourcesClient) BackupInstancesDeleteThenPoll(ctx context.Context, id BackupInstanceId, options BackupInstancesDeleteOperationOptions) error {
	result, err := c.BackupInstancesDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing BackupInstancesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackupInstancesDelete: %+v", err)
	}

	return nil
}
