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

type TableServicesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListTableServices
}

// TableServicesList ...
func (c TableServicePropertiesClient) TableServicesList(ctx context.Context, id commonids.StorageAccountId) (result TableServicesListOperationResponse, err error) {
	req, err := c.preparerForTableServicesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableserviceproperties.TableServicePropertiesClient", "TableServicesList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableserviceproperties.TableServicePropertiesClient", "TableServicesList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableServicesList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableserviceproperties.TableServicePropertiesClient", "TableServicesList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableServicesList prepares the TableServicesList request.
func (c TableServicePropertiesClient) preparerForTableServicesList(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tableServices", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableServicesList handles the response to the TableServicesList request. The method always
// closes the http.Response Body.
func (c TableServicePropertiesClient) responderForTableServicesList(resp *http.Response) (result TableServicesListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
