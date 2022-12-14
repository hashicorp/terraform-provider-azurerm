package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainersListLogsOperationResponse struct {
	HttpResponse *http.Response
	Model        *Logs
}

type ContainersListLogsOperationOptions struct {
	Tail       *int64
	Timestamps *bool
}

func DefaultContainersListLogsOperationOptions() ContainersListLogsOperationOptions {
	return ContainersListLogsOperationOptions{}
}

func (o ContainersListLogsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ContainersListLogsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Tail != nil {
		out["tail"] = *o.Tail
	}

	if o.Timestamps != nil {
		out["timestamps"] = *o.Timestamps
	}

	return out
}

// ContainersListLogs ...
func (c ContainerInstanceClient) ContainersListLogs(ctx context.Context, id ContainerId, options ContainersListLogsOperationOptions) (result ContainersListLogsOperationResponse, err error) {
	req, err := c.preparerForContainersListLogs(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersListLogs", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersListLogs", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainersListLogs(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersListLogs", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainersListLogs prepares the ContainersListLogs request.
func (c ContainerInstanceClient) preparerForContainersListLogs(ctx context.Context, id ContainerId, options ContainersListLogsOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/logs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainersListLogs handles the response to the ContainersListLogs request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainersListLogs(resp *http.Response) (result ContainersListLogsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
