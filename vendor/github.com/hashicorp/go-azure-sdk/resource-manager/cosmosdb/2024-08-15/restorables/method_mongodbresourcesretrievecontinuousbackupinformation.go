package restorables

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

type MongoDBResourcesRetrieveContinuousBackupInformationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BackupInformation
}

// MongoDBResourcesRetrieveContinuousBackupInformation ...
func (c RestorablesClient) MongoDBResourcesRetrieveContinuousBackupInformation(ctx context.Context, id MongodbDatabaseCollectionId, input ContinuousBackupRestoreLocation) (result MongoDBResourcesRetrieveContinuousBackupInformationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/retrieveContinuousBackupInformation", id.ID()),
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

// MongoDBResourcesRetrieveContinuousBackupInformationThenPoll performs MongoDBResourcesRetrieveContinuousBackupInformation then polls until it's completed
func (c RestorablesClient) MongoDBResourcesRetrieveContinuousBackupInformationThenPoll(ctx context.Context, id MongodbDatabaseCollectionId, input ContinuousBackupRestoreLocation) error {
	result, err := c.MongoDBResourcesRetrieveContinuousBackupInformation(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesRetrieveContinuousBackupInformation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesRetrieveContinuousBackupInformation: %+v", err)
	}

	return nil
}
