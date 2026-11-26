package vpnserverconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RadiusAuthServerOperationPredicate struct {
	RadiusServerAddress *string
	RadiusServerSecret  *string
}

func (p RadiusAuthServerOperationPredicate) Matches(input RadiusAuthServer) bool {

	if p.RadiusServerAddress != nil && (input.RadiusServerAddress == nil || *p.RadiusServerAddress != *input.RadiusServerAddress) {
		return false
	}

	if p.RadiusServerSecret != nil && (input.RadiusServerSecret == nil || *p.RadiusServerSecret != *input.RadiusServerSecret) {
		return false
	}

	return true
}
