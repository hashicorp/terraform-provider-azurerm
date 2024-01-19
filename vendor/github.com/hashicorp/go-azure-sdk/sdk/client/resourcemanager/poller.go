// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

func PollerFromResponse(response *client.Response, client *Client) (poller pollers.Poller, err error) {
	if response == nil {
		return pollers.Poller{}, fmt.Errorf("no HTTP Response was returned")
	}

	originalRequestUri := response.Request.URL.String()

	// If this is a LRO we should either have a 201/202 with a Polling URI header
	isLroStatus := response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusAccepted
	methodIsDelete := strings.EqualFold(response.Request.Method, "DELETE")
	lroPollingUri := pollingUriForLongRunningOperation(response)
	lroIsSelfReference := isLROSelfReference(lroPollingUri, originalRequestUri)
	if isLroStatus && lroPollingUri != "" && !methodIsDelete && !lroIsSelfReference {
		lro, lroErr := longRunningOperationPollerFromResponse(response, client.Client)
		if lroErr != nil {
			err = lroErr
			return pollers.Poller{}, fmt.Errorf("building long-running-operation poller: %+v", lroErr)
		}
		return pollers.NewPoller(lro, lro.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
	}

	// or we should be polling on the `provisioningState` of the resource
	contentType := response.Header.Get("Content-Type")
	if contentType == "" && response.Request != nil {
		contentType = response.Request.Header.Get("Content-Type")
	}

	statusCodesToCheckProvisioningState := response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated || (response.StatusCode == http.StatusAccepted && lroIsSelfReference)
	contentTypeMatchesForProvisioningStateCheck := strings.Contains(strings.ToLower(contentType), "application/json")
	methodIsApplicable := strings.EqualFold(response.Request.Method, "PATCH") ||
		strings.EqualFold(response.Request.Method, "POST") ||
		strings.EqualFold(response.Request.Method, "PUT")
	if statusCodesToCheckProvisioningState && contentTypeMatchesForProvisioningStateCheck && methodIsApplicable {
		provisioningState, provisioningStateErr := provisioningStatePollerFromResponse(response, lroIsSelfReference, client, DefaultPollingInterval)
		if provisioningStateErr != nil {
			err = provisioningStateErr
			return pollers.Poller{}, fmt.Errorf("building provisioningState poller: %+v", provisioningStateErr)
		}
		return pollers.NewPoller(provisioningState, provisioningState.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
	}

	statusCodesToCheckDelete := response.StatusCode == http.StatusOK || response.StatusCode == http.StatusNoContent
	statusCodesToCheckLroDelete := response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusAccepted
	if methodIsDelete {
		// finally, if it was a Delete that returned a 200/204
		if statusCodesToCheckDelete {
			deletePoller, deletePollerErr := deletePollerFromResponse(response, client, DefaultPollingInterval)
			if deletePollerErr != nil {
				err = deletePollerErr
				return pollers.Poller{}, fmt.Errorf("building delete poller: %+v", deletePollerErr)
			}
			return pollers.NewPoller(deletePoller, deletePoller.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
		}else {
			// finally, if it was a Delete that returned a 201/202
			// For APIM servcie, even if Get API returns 404, deletion is still in progress.
			// This leads Terraform to believe that the deletion has been completed but is actually still being deleted.
			// This will cause the resource group deletion to fail because the APIM servcie still exists.
			// Feedback from the service team that
			// for long-running-operation deletion, track asynchronous Azure operations instead of whether the get api returns 404 can more accurately determine whether the operation is actually completed.
			if statusCodesToCheckLroDelete{
				lro, lroErr := longRunningOperationPollerFromResponse(response, client.Client)
				if lroErr != nil {
					err = lroErr
					return pollers.Poller{}, fmt.Errorf("building long-running-operation poller: %+v", lroErr)
				}
				return pollers.NewPoller(lro, lro.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
			}
		}
	}

	return pollers.Poller{}, fmt.Errorf("no applicable pollers were found for the response")
}

func isLROSelfReference(lroPollingUri, originalRequestUri string) bool {
	// Some APIs return a LRO URI of themselves, meaning that we should be checking a 200 OK is returned rather
	// than polling as usual. Automation@2022-08-08 - DSCNodeConfiguration CreateOrUpdate is one such example.
	first, err := url.Parse(lroPollingUri)
	if err != nil {
		return false
	}
	second, err := url.Parse(originalRequestUri)
	if err != nil {
		return false
	}

	// The Query String can be a different API version / options and in some cases the Host
	// is returned with `:443` - so the path should be sufficient as a check.
	return strings.EqualFold(first.Path, second.Path)
}
