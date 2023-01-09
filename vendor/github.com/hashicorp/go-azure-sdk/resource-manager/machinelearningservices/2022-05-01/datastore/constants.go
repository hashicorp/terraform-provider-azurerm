package datastore

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialsType string

const (
	CredentialsTypeAccountKey       CredentialsType = "AccountKey"
	CredentialsTypeCertificate      CredentialsType = "Certificate"
	CredentialsTypeNone             CredentialsType = "None"
	CredentialsTypeSas              CredentialsType = "Sas"
	CredentialsTypeServicePrincipal CredentialsType = "ServicePrincipal"
)

func PossibleValuesForCredentialsType() []string {
	return []string{
		string(CredentialsTypeAccountKey),
		string(CredentialsTypeCertificate),
		string(CredentialsTypeNone),
		string(CredentialsTypeSas),
		string(CredentialsTypeServicePrincipal),
	}
}

func parseCredentialsType(input string) (*CredentialsType, error) {
	vals := map[string]CredentialsType{
		"accountkey":       CredentialsTypeAccountKey,
		"certificate":      CredentialsTypeCertificate,
		"none":             CredentialsTypeNone,
		"sas":              CredentialsTypeSas,
		"serviceprincipal": CredentialsTypeServicePrincipal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CredentialsType(input)
	return &out, nil
}

type DatastoreType string

const (
	DatastoreTypeAzureBlob           DatastoreType = "AzureBlob"
	DatastoreTypeAzureDataLakeGenOne DatastoreType = "AzureDataLakeGen1"
	DatastoreTypeAzureDataLakeGenTwo DatastoreType = "AzureDataLakeGen2"
	DatastoreTypeAzureFile           DatastoreType = "AzureFile"
)

func PossibleValuesForDatastoreType() []string {
	return []string{
		string(DatastoreTypeAzureBlob),
		string(DatastoreTypeAzureDataLakeGenOne),
		string(DatastoreTypeAzureDataLakeGenTwo),
		string(DatastoreTypeAzureFile),
	}
}

func parseDatastoreType(input string) (*DatastoreType, error) {
	vals := map[string]DatastoreType{
		"azureblob":         DatastoreTypeAzureBlob,
		"azuredatalakegen1": DatastoreTypeAzureDataLakeGenOne,
		"azuredatalakegen2": DatastoreTypeAzureDataLakeGenTwo,
		"azurefile":         DatastoreTypeAzureFile,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatastoreType(input)
	return &out, nil
}

type SecretsType string

const (
	SecretsTypeAccountKey       SecretsType = "AccountKey"
	SecretsTypeCertificate      SecretsType = "Certificate"
	SecretsTypeSas              SecretsType = "Sas"
	SecretsTypeServicePrincipal SecretsType = "ServicePrincipal"
)

func PossibleValuesForSecretsType() []string {
	return []string{
		string(SecretsTypeAccountKey),
		string(SecretsTypeCertificate),
		string(SecretsTypeSas),
		string(SecretsTypeServicePrincipal),
	}
}

func parseSecretsType(input string) (*SecretsType, error) {
	vals := map[string]SecretsType{
		"accountkey":       SecretsTypeAccountKey,
		"certificate":      SecretsTypeCertificate,
		"sas":              SecretsTypeSas,
		"serviceprincipal": SecretsTypeServicePrincipal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecretsType(input)
	return &out, nil
}

type ServiceDataAccessAuthIdentity string

const (
	ServiceDataAccessAuthIdentityNone                            ServiceDataAccessAuthIdentity = "None"
	ServiceDataAccessAuthIdentityWorkspaceSystemAssignedIdentity ServiceDataAccessAuthIdentity = "WorkspaceSystemAssignedIdentity"
	ServiceDataAccessAuthIdentityWorkspaceUserAssignedIdentity   ServiceDataAccessAuthIdentity = "WorkspaceUserAssignedIdentity"
)

func PossibleValuesForServiceDataAccessAuthIdentity() []string {
	return []string{
		string(ServiceDataAccessAuthIdentityNone),
		string(ServiceDataAccessAuthIdentityWorkspaceSystemAssignedIdentity),
		string(ServiceDataAccessAuthIdentityWorkspaceUserAssignedIdentity),
	}
}

func parseServiceDataAccessAuthIdentity(input string) (*ServiceDataAccessAuthIdentity, error) {
	vals := map[string]ServiceDataAccessAuthIdentity{
		"none":                            ServiceDataAccessAuthIdentityNone,
		"workspacesystemassignedidentity": ServiceDataAccessAuthIdentityWorkspaceSystemAssignedIdentity,
		"workspaceuserassignedidentity":   ServiceDataAccessAuthIdentityWorkspaceUserAssignedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceDataAccessAuthIdentity(input)
	return &out, nil
}
