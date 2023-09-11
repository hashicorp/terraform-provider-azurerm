package assessmentsmetadata

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentsMetadataGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SecurityAssessmentMetadataResponse
}

// AssessmentsMetadataGet ...
func (c AssessmentsMetadataClient) AssessmentsMetadataGet(ctx context.Context, id AssessmentMetadataId) (result AssessmentsMetadataGetOperationResponse, err error) {
	req, err := c.preparerForAssessmentsMetadataGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssessmentsMetadataGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssessmentsMetadataGet prepares the AssessmentsMetadataGet request.
func (c AssessmentsMetadataClient) preparerForAssessmentsMetadataGet(ctx context.Context, id AssessmentMetadataId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAssessmentsMetadataGet handles the response to the AssessmentsMetadataGet request. The method always
// closes the http.Response Body.
func (c AssessmentsMetadataClient) responderForAssessmentsMetadataGet(resp *http.Response) (result AssessmentsMetadataGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
