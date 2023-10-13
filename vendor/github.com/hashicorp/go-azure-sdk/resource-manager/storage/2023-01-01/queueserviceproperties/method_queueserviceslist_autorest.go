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

type QueueServicesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListQueueServices
}

// QueueServicesList ...
func (c QueueServicePropertiesClient) QueueServicesList(ctx context.Context, id commonids.StorageAccountId) (result QueueServicesListOperationResponse, err error) {
	req, err := c.preparerForQueueServicesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueServicesList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueserviceproperties.QueueServicePropertiesClient", "QueueServicesList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueServicesList prepares the QueueServicesList request.
func (c QueueServicePropertiesClient) preparerForQueueServicesList(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/queueServices", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForQueueServicesList handles the response to the QueueServicesList request. The method always
// closes the http.Response Body.
func (c QueueServicePropertiesClient) responderForQueueServicesList(resp *http.Response) (result QueueServicesListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
