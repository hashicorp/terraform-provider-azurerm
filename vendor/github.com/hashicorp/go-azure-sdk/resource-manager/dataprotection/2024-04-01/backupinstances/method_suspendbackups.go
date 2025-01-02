package backupinstances

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

type SuspendBackupsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type SuspendBackupsOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultSuspendBackupsOperationOptions() SuspendBackupsOperationOptions {
	return SuspendBackupsOperationOptions{}
}

func (o SuspendBackupsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsAuthorizationAuxiliary != nil {
		out.Append("x-ms-authorization-auxiliary", fmt.Sprintf("%v", *o.XMsAuthorizationAuxiliary))
	}
	return &out
}

func (o SuspendBackupsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SuspendBackupsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// SuspendBackups ...
func (c BackupInstancesClient) SuspendBackups(ctx context.Context, id BackupInstanceId, input SuspendBackupRequest, options SuspendBackupsOperationOptions) (result SuspendBackupsOperationResponse, err error) {
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

// SuspendBackupsThenPoll performs SuspendBackups then polls until it's completed
func (c BackupInstancesClient) SuspendBackupsThenPoll(ctx context.Context, id BackupInstanceId, input SuspendBackupRequest, options SuspendBackupsOperationOptions) error {
	result, err := c.SuspendBackups(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing SuspendBackups: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SuspendBackups: %+v", err)
	}

	return nil
}
