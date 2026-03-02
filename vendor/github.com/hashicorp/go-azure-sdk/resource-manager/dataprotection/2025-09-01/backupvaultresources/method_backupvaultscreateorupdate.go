package backupvaultresources

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

type BackupVaultsCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BackupVaultResource
}

type BackupVaultsCreateOrUpdateOperationOptions struct {
	XMsAuthorizationAuxiliary *string
	XMsDeletedVaultId         *string
}

func DefaultBackupVaultsCreateOrUpdateOperationOptions() BackupVaultsCreateOrUpdateOperationOptions {
	return BackupVaultsCreateOrUpdateOperationOptions{}
}

func (o BackupVaultsCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsAuthorizationAuxiliary != nil {
		out.Append("x-ms-authorization-auxiliary", fmt.Sprintf("%v", *o.XMsAuthorizationAuxiliary))
	}
	if o.XMsDeletedVaultId != nil {
		out.Append("x-ms-deleted-vault-id", fmt.Sprintf("%v", *o.XMsDeletedVaultId))
	}
	return &out
}

func (o BackupVaultsCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o BackupVaultsCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// BackupVaultsCreateOrUpdate ...
func (c BackupVaultResourcesClient) BackupVaultsCreateOrUpdate(ctx context.Context, id BackupVaultId, input BackupVaultResource, options BackupVaultsCreateOrUpdateOperationOptions) (result BackupVaultsCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
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

// BackupVaultsCreateOrUpdateThenPoll performs BackupVaultsCreateOrUpdate then polls until it's completed
func (c BackupVaultResourcesClient) BackupVaultsCreateOrUpdateThenPoll(ctx context.Context, id BackupVaultId, input BackupVaultResource, options BackupVaultsCreateOrUpdateOperationOptions) error {
	result, err := c.BackupVaultsCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing BackupVaultsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackupVaultsCreateOrUpdate: %+v", err)
	}

	return nil
}
