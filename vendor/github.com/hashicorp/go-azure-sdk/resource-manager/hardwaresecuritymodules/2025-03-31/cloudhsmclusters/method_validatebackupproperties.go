package cloudhsmclusters

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

type ValidateBackupPropertiesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BackupResult
}

// ValidateBackupProperties ...
func (c CloudHsmClustersClient) ValidateBackupProperties(ctx context.Context, id CloudHsmClusterId, input BackupRestoreRequestBaseProperties) (result ValidateBackupPropertiesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/validateBackupProperties", id.ID()),
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

// ValidateBackupPropertiesThenPoll performs ValidateBackupProperties then polls until it's completed
func (c CloudHsmClustersClient) ValidateBackupPropertiesThenPoll(ctx context.Context, id CloudHsmClusterId, input BackupRestoreRequestBaseProperties) error {
	result, err := c.ValidateBackupProperties(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ValidateBackupProperties: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ValidateBackupProperties: %+v", err)
	}

	return nil
}
