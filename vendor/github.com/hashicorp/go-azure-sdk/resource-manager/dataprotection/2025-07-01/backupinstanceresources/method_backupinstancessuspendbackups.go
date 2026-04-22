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

type BackupInstancesSuspendBackupsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type BackupInstancesSuspendBackupsOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultBackupInstancesSuspendBackupsOperationOptions() BackupInstancesSuspendBackupsOperationOptions {
	return BackupInstancesSuspendBackupsOperationOptions{}
}

func (o BackupInstancesSuspendBackupsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsAuthorizationAuxiliary != nil {
		out.Append("x-ms-authorization-auxiliary", fmt.Sprintf("%v", *o.XMsAuthorizationAuxiliary))
	}
	return &out
}

func (o BackupInstancesSuspendBackupsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o BackupInstancesSuspendBackupsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// BackupInstancesSuspendBackups ...
func (c BackupInstanceResourcesClient) BackupInstancesSuspendBackups(ctx context.Context, id BackupInstanceId, input SuspendBackupRequest, options BackupInstancesSuspendBackupsOperationOptions) (result BackupInstancesSuspendBackupsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/suspendBackups", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

// BackupInstancesSuspendBackupsThenPoll performs BackupInstancesSuspendBackups then polls until it's completed
func (c BackupInstanceResourcesClient) BackupInstancesSuspendBackupsThenPoll(ctx context.Context, id BackupInstanceId, input SuspendBackupRequest, options BackupInstancesSuspendBackupsOperationOptions) error {
	result, err := c.BackupInstancesSuspendBackups(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing BackupInstancesSuspendBackups: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackupInstancesSuspendBackups: %+v", err)
	}

	return nil
}
