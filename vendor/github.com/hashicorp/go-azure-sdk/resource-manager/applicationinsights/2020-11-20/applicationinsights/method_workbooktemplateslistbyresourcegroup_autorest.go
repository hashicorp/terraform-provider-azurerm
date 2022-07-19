package applicationinsights

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

type WorkbookTemplatesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkbookTemplatesListResult
}

// WorkbookTemplatesListByResourceGroup ...
func (c ApplicationInsightsClient) WorkbookTemplatesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result WorkbookTemplatesListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForWorkbookTemplatesListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbookTemplatesListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbookTemplatesListByResourceGroup prepares the WorkbookTemplatesListByResourceGroup request.
func (c ApplicationInsightsClient) preparerForWorkbookTemplatesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/workbookTemplates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbookTemplatesListByResourceGroup handles the response to the WorkbookTemplatesListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbookTemplatesListByResourceGroup(resp *http.Response) (result WorkbookTemplatesListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
