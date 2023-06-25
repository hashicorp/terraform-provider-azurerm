package projectresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

type ProjectsDeleteOperationOptions struct {
	DeleteRunningTasks *bool
}

func DefaultProjectsDeleteOperationOptions() ProjectsDeleteOperationOptions {
	return ProjectsDeleteOperationOptions{}
}

func (o ProjectsDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ProjectsDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.DeleteRunningTasks != nil {
		out["deleteRunningTasks"] = *o.DeleteRunningTasks
	}

	return out
}

// ProjectsDelete ...
func (c ProjectResourceClient) ProjectsDelete(ctx context.Context, id ProjectId, options ProjectsDeleteOperationOptions) (result ProjectsDeleteOperationResponse, err error) {
	req, err := c.preparerForProjectsDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForProjectsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForProjectsDelete prepares the ProjectsDelete request.
func (c ProjectResourceClient) preparerForProjectsDelete(ctx context.Context, id ProjectId, options ProjectsDeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForProjectsDelete handles the response to the ProjectsDelete request. The method always
// closes the http.Response Body.
func (c ProjectResourceClient) responderForProjectsDelete(resp *http.Response) (result ProjectsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
