package nginxconfiguration

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

type ConfigurationsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConfigurationsDelete ...
func (c NginxConfigurationClient) ConfigurationsDelete(ctx context.Context, id ConfigurationId) (result ConfigurationsDeleteOperationResponse, err error) {
	req, err := c.preparerForConfigurationsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxconfiguration.NginxConfigurationClient", "ConfigurationsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConfigurationsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxconfiguration.NginxConfigurationClient", "ConfigurationsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConfigurationsDeleteThenPoll performs ConfigurationsDelete then polls until it's completed
func (c NginxConfigurationClient) ConfigurationsDeleteThenPoll(ctx context.Context, id ConfigurationId) error {
	result, err := c.ConfigurationsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ConfigurationsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConfigurationsDelete: %+v", err)
	}

	return nil
}

// preparerForConfigurationsDelete prepares the ConfigurationsDelete request.
func (c NginxConfigurationClient) preparerForConfigurationsDelete(ctx context.Context, id ConfigurationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForConfigurationsDelete sends the ConfigurationsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c NginxConfigurationClient) senderForConfigurationsDelete(ctx context.Context, req *http.Request) (future ConfigurationsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
