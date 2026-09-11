package expressrouteserviceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteServiceProviderBandwidthsOffered struct {
	OfferName   *string `json:"offerName,omitempty"`
	ValueInMbps *int64  `json:"valueInMbps,omitempty"`
}
