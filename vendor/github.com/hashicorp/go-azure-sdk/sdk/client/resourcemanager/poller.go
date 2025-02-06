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

	// If this is a LRO we should either have a 200/201/202 with a Polling URI header
	isLroStatus := response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusAccepted
	methodIsDelete := strings.EqualFold(response.Request.Method, "DELETE")
	lroPollingUri := pollingUriForLongRunningOperation(response)
	lroIsSelfReference := isLROSelfReference(lroPollingUri, originalRequestUri)
	if isLroStatus && lroPollingUri != "" && !methodIsDelete && !lroIsSelfReference {
		lro, lroErr := longRunningOperationPollerFromResponse(response, client.Client)
		if lroErr != nil {
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
			return pollers.Poller{}, fmt.Errorf("building provisioningState poller: %+v", provisioningStateErr)
		}
		return pollers.NewPoller(provisioningState, provisioningState.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
	}

	// finally, if it was a Delete that returned a 200/204
	statusCodesToCheckDelete := response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusAccepted || response.StatusCode == http.StatusNoContent
	if methodIsDelete && statusCodesToCheckDelete {
		deletePoller, deletePollerErr := deletePollerFromResponse(response, client, DefaultPollingInterval)
		if deletePollerErr != nil {
			return pollers.Poller{}, fmt.Errorf("building delete poller: %+v", deletePollerErr)
		}
		return pollers.NewPoller(deletePoller, deletePoller.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
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
