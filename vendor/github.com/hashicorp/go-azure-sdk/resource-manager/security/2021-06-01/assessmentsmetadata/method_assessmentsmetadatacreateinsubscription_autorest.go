package assessmentsmetadata

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentsMetadataCreateInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *SecurityAssessmentMetadataResponse
}

// AssessmentsMetadataCreateInSubscription ...
func (c AssessmentsMetadataClient) AssessmentsMetadataCreateInSubscription(ctx context.Context, id ProviderAssessmentMetadataId, input SecurityAssessmentMetadataResponse) (result AssessmentsMetadataCreateInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForAssessmentsMetadataCreateInSubscription(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataCreateInSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataCreateInSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssessmentsMetadataCreateInSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataCreateInSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssessmentsMetadataCreateInSubscription prepares the AssessmentsMetadataCreateInSubscription request.
func (c AssessmentsMetadataClient) preparerForAssessmentsMetadataCreateInSubscription(ctx context.Context, id ProviderAssessmentMetadataId, input SecurityAssessmentMetadataResponse) (*http.Request, error) {
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

// responderForAssessmentsMetadataCreateInSubscription handles the response to the AssessmentsMetadataCreateInSubscription request. The method always
// closes the http.Response Body.
func (c AssessmentsMetadataClient) responderForAssessmentsMetadataCreateInSubscription(resp *http.Response) (result AssessmentsMetadataCreateInSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
