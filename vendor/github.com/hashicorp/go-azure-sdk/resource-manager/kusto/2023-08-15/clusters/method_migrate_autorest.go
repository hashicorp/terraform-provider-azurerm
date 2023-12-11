package clusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Migrate ...
func (c ClustersClient) Migrate(ctx context.Context, id ClusterId, input ClusterMigrateRequest) (result MigrateOperationResponse, err error) {
	req, err := c.preparerForMigrate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "Migrate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMigrate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "Migrate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MigrateThenPoll performs Migrate then polls until it's completed
func (c ClustersClient) MigrateThenPoll(ctx context.Context, id ClusterId, input ClusterMigrateRequest) error {
	result, err := c.Migrate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Migrate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Migrate: %+v", err)
	}

	return nil
}

// preparerForMigrate prepares the Migrate request.
func (c ClustersClient) preparerForMigrate(ctx context.Context, id ClusterId, input ClusterMigrateRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/migrate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMigrate sends the Migrate request. The method will close the
// http.Response Body if it receives an error.
func (c ClustersClient) senderForMigrate(ctx context.Context, req *http.Request) (future MigrateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
