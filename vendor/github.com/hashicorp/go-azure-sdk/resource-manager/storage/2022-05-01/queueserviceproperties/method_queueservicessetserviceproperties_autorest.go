package queueserviceproperties

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueServicesSetServicePropertiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *QueueServiceProperties
}

// QueueServicesSetServiceProperties ...
func (c QueueServicePropertiesClient) QueueServicesSetServiceProperties(ctx context.Context, id commonids.StorageAccountId, input QueueServiceProperties) (result QueueServicesSetServicePropertiesOperationResponse, err error) {
	req, err := c.preparerForQueueServicesSetServiceProperties(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesSetServiceProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesSetServiceProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueServicesSetServiceProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesSetServiceProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueServicesSetServiceProperties prepares the QueueServicesSetServiceProperties request.
func (c QueueServicePropertiesClient) preparerForQueueServicesSetServiceProperties(ctx context.Context, id commonids.StorageAccountId, input QueueServiceProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/queueServices/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForQueueServicesSetServiceProperties handles the response to the QueueServicesSetServiceProperties request. The method always
// closes the http.Response Body.
func (c QueueServicePropertiesClient) responderForQueueServicesSetServiceProperties(resp *http.Response) (result QueueServicesSetServicePropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
