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

type QueueServicesGetServicePropertiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *QueueServiceProperties
}

// QueueServicesGetServiceProperties ...
func (c QueueServicePropertiesClient) QueueServicesGetServiceProperties(ctx context.Context, id commonids.StorageAccountId) (result QueueServicesGetServicePropertiesOperationResponse, err error) {
	req, err := c.preparerForQueueServicesGetServiceProperties(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesGetServiceProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesGetServiceProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueServicesGetServiceProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesGetServiceProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueServicesGetServiceProperties prepares the QueueServicesGetServiceProperties request.
func (c QueueServicePropertiesClient) preparerForQueueServicesGetServiceProperties(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/queueServices/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForQueueServicesGetServiceProperties handles the response to the QueueServicesGetServiceProperties request. The method always
// closes the http.Response Body.
func (c QueueServicePropertiesClient) responderForQueueServicesGetServiceProperties(resp *http.Response) (result QueueServicesGetServicePropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
