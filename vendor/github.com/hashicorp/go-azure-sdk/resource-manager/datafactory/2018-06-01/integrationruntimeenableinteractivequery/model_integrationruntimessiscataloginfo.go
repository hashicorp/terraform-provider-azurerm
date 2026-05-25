package integrationruntimeenableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeSsisCatalogInfo struct {
	CatalogAdminPassword  *SecureString                             `json:"catalogAdminPassword,omitempty"`
	CatalogAdminUserName  *string                                   `json:"catalogAdminUserName,omitempty"`
	CatalogPricingTier    *IntegrationRuntimeSsisCatalogPricingTier `json:"catalogPricingTier,omitempty"`
	CatalogServerEndpoint *string                                   `json:"catalogServerEndpoint,omitempty"`
	DualStandbyPairName   *string                                   `json:"dualStandbyPairName,omitempty"`
}
