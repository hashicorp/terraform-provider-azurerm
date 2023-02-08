package singlesignon

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

type ConfigurationsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConfigurationsCreateOrUpdate ...
func (c SingleSignOnClient) ConfigurationsCreateOrUpdate(ctx context.Context, id SingleSignOnConfigurationId, input DatadogSingleSignOnResource) (result ConfigurationsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForConfigurationsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "singlesignon.SingleSignOnClient", "ConfigurationsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConfigurationsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "singlesignon.SingleSignOnClient", "ConfigurationsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConfigurationsCreateOrUpdateThenPoll performs ConfigurationsCreateOrUpdate then polls until it's completed
func (c SingleSignOnClient) ConfigurationsCreateOrUpdateThenPoll(ctx context.Context, id SingleSignOnConfigurationId, input DatadogSingleSignOnResource) error {
	result, err := c.ConfigurationsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ConfigurationsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConfigurationsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForConfigurationsCreateOrUpdate prepares the ConfigurationsCreateOrUpdate request.
func (c SingleSignOnClient) preparerForConfigurationsCreateOrUpdate(ctx context.Context, id SingleSignOnConfigurationId, input DatadogSingleSignOnResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForConfigurationsCreateOrUpdate sends the ConfigurationsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c SingleSignOnClient) senderForConfigurationsCreateOrUpdate(ctx context.Context, req *http.Request) (future ConfigurationsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
