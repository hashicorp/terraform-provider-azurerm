// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metadata

type MetaData struct {
	Authentication          Authentication
	DnsSuffixes             DnsSuffixes
	Name                    string
	ResourceIdentifiers     ResourceIdentifiers
	ResourceManagerEndpoint string
}

type Authentication struct {
	Audiences        []string
	LoginEndpoint    string
	IdentityProvider string
	Tenant           string
}

type DnsSuffixes struct {
	Attestation string
	FrontDoor   string
	KeyVault    string
	ManagedHSM  string
	MariaDB     string
	MySql       string
	Postgresql  string
	SqlServer   string
	Storage     string
	StorageSync string
	Synapse     string
}

type ResourceIdentifiers struct {
	Attestation    string
	Batch          string
	LogAnalytics   string
	Media          string
	MicrosoftGraph string
	OSSRDBMS       string
	Synapse        string
}
