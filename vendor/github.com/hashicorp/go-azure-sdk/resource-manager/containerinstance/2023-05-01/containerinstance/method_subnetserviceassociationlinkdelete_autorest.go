package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubnetServiceAssociationLinkDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SubnetServiceAssociationLinkDelete ...
func (c ContainerInstanceClient) SubnetServiceAssociationLinkDelete(ctx context.Context, id commonids.SubnetId) (result SubnetServiceAssociationLinkDeleteOperationResponse, err error) {
	req, err := c.preparerForSubnetServiceAssociationLinkDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "SubnetServiceAssociationLinkDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSubnetServiceAssociationLinkDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "SubnetServiceAssociationLinkDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SubnetServiceAssociationLinkDeleteThenPoll performs SubnetServiceAssociationLinkDelete then polls until it's completed
func (c ContainerInstanceClient) SubnetServiceAssociationLinkDeleteThenPoll(ctx context.Context, id commonids.SubnetId) error {
	result, err := c.SubnetServiceAssociationLinkDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SubnetServiceAssociationLinkDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SubnetServiceAssociationLinkDelete: %+v", err)
	}

	return nil
}

// preparerForSubnetServiceAssociationLinkDelete prepares the SubnetServiceAssociationLinkDelete request.
func (c ContainerInstanceClient) preparerForSubnetServiceAssociationLinkDelete(ctx context.Context, id commonids.SubnetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ContainerInstance/serviceAssociationLinks/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSubnetServiceAssociationLinkDelete sends the SubnetServiceAssociationLinkDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ContainerInstanceClient) senderForSubnetServiceAssociationLinkDelete(ctx context.Context, req *http.Request) (future SubnetServiceAssociationLinkDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
