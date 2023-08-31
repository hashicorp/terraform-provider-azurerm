package tableserviceproperties

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

type TableServicesGetServicePropertiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *TableServiceProperties
}

// TableServicesGetServiceProperties ...
func (c TableServicePropertiesClient) TableServicesGetServiceProperties(ctx context.Context, id commonids.StorageAccountId) (result TableServicesGetServicePropertiesOperationResponse, err error) {
	req, err := c.preparerForTableServicesGetServiceProperties(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableserviceproperties.TableServicePropertiesClient", "TableServicesGetServiceProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableserviceproperties.TableServicePropertiesClient", "TableServicesGetServiceProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableServicesGetServiceProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableserviceproperties.TableServicePropertiesClient", "TableServicesGetServiceProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableServicesGetServiceProperties prepares the TableServicesGetServiceProperties request.
func (c TableServicePropertiesClient) preparerForTableServicesGetServiceProperties(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tableServices/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableServicesGetServiceProperties handles the response to the TableServicesGetServiceProperties request. The method always
// closes the http.Response Body.
func (c TableServicePropertiesClient) responderForTableServicesGetServiceProperties(resp *http.Response) (result TableServicesGetServicePropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
