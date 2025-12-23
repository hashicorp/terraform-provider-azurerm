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

type BackupInstancesTriggerRestoreOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *OperationJobExtendedInfo
}

type BackupInstancesTriggerRestoreOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultBackupInstancesTriggerRestoreOperationOptions() BackupInstancesTriggerRestoreOperationOptions {
	return BackupInstancesTriggerRestoreOperationOptions{}
}

func (o BackupInstancesTriggerRestoreOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsAuthorizationAuxiliary != nil {
		out.Append("x-ms-authorization-auxiliary", fmt.Sprintf("%v", *o.XMsAuthorizationAuxiliary))
	}
	return &out
}

func (o BackupInstancesTriggerRestoreOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o BackupInstancesTriggerRestoreOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// BackupInstancesTriggerRestore ...
func (c BackupInstanceResourcesClient) BackupInstancesTriggerRestore(ctx context.Context, id BackupInstanceId, input AzureBackupRestoreRequest, options BackupInstancesTriggerRestoreOperationOptions) (result BackupInstancesTriggerRestoreOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restore", id.ID()),
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

// BackupInstancesTriggerRestoreThenPoll performs BackupInstancesTriggerRestore then polls until it's completed
func (c BackupInstanceResourcesClient) BackupInstancesTriggerRestoreThenPoll(ctx context.Context, id BackupInstanceId, input AzureBackupRestoreRequest, options BackupInstancesTriggerRestoreOperationOptions) error {
	result, err := c.BackupInstancesTriggerRestore(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing BackupInstancesTriggerRestore: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackupInstancesTriggerRestore: %+v", err)
	}

	return nil
}
