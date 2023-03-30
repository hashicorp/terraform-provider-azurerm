// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/internal/stringfmt"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// TODO: return a typed error here so that we could potentially change this error/expose this in the Provider
// TODO: the error should return the default error message shown below

var defaultRetryFunctions = []client.RequestRetryFunc{
	// NOTE: 429 is handled by the base library
	handleResourceProviderNotRegistered,
}

func handleResourceProviderNotRegistered(r *http.Response, o *odata.OData) (bool, error) {
	if o != nil && o.Error != nil && o.Error.Code != nil && strings.EqualFold(*o.Error.Code, "MissingSubscriptionRegistration") {
		return false, resourceProviderNotRegisteredError(*o.Error.Message)
	}

	return false, nil
}

func resourceProviderNotRegisteredError(message string) error {
	messageSplit := stringfmt.QuoteAndSplitString(">", message, 100)
	return fmt.Errorf(`The Resource Provider was not registered

Resource Providers (APIs) in Azure need to be registered before they can be used - however the Resource
Provider was not registered, and calling the API returned the following error:

%[1]s

The Azure Provider by default will automatically register certain Resource Providers at launch-time,
whilst it's possible to opt-out of this (which you may have done) 

Please ensure that this Resource Provider is properly registered, you can do this using the Azure CLI
for example to register the Resource Provider "Some.ResourceProvider" is registered run:

> az provider register --namespace "Some.ResourceProvider"

Resource Providers can take a while to register, you can check the status by running:

> az provider show --namespace "Some.ResourceProvider" --query "registrationState"

Once this outputs "Registered" the Resource Provider is available for use and you can re-run Terraform.
`, strings.Join(messageSplit, "\n"))
}
