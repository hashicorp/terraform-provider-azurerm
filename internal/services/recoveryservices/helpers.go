// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/lang/response"
)

// This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
func handleAzureSdkForGoBug2824(id string) string {
	return strings.Replace(id, "/Subscriptions/", "/subscriptions/", 1)
}

func wasBadRequestWithNotExist(resp *http.Response, err error) bool {
	e, ok := err.(autorest.DetailedError)
	if !ok {
		return false
	}

	r, ok := e.Original.(*azure.RequestError)
	if !ok {
		return false
	}

	if r.ServiceError == nil || len(r.ServiceError.Details) == 0 {
		return false
	}

	sc, ok := r.ServiceError.Details[0]["code"]
	if !ok {
		return false
	}

	return response.WasBadRequest(resp) && sc == "SubscriptionIdNotRegisteredWithSrs"
}
