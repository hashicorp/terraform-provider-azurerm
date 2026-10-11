package v7_4

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedcertificates"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedkeys"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedsecrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedstorage"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/fullbackup"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/fullrestore"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/hsmsecuritydomain"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/rng"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/roleassignments"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/roledefinitions"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/settings"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/storage"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

type Client struct {
	Certificates        *certificates.CertificatesClient
	DeletedCertificates *deletedcertificates.DeletedCertificatesClient
	DeletedKeys         *deletedkeys.DeletedKeysClient
	DeletedSecrets      *deletedsecrets.DeletedSecretsClient
	DeletedStorage      *deletedstorage.DeletedStorageClient
	FullBackup          *fullbackup.FullBackupClient
	FullRestore         *fullrestore.FullRestoreClient
	HSMSecurityDomain   *hsmsecuritydomain.HSMSecurityDomainClient
	Keys                *keys.KeysClient
	RNG                 *rng.RNGClient
	RoleAssignments     *roleassignments.RoleAssignmentsClient
	RoleDefinitions     *roledefinitions.RoleDefinitionsClient
	Secrets             *secrets.SecretsClient
	Settings            *settings.SettingsClient
	Storage             *storage.StorageClient
}

func NewClient(configureFunc func(c *dataplane.Client)) (*Client, error) {
	certificatesClient, err := certificates.NewCertificatesClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Certificates client: %+v", err)
	}
	configureFunc(certificatesClient.Client)

	deletedCertificatesClient, err := deletedcertificates.NewDeletedCertificatesClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building DeletedCertificates client: %+v", err)
	}
	configureFunc(deletedCertificatesClient.Client)

	deletedKeysClient, err := deletedkeys.NewDeletedKeysClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building DeletedKeys client: %+v", err)
	}
	configureFunc(deletedKeysClient.Client)

	deletedSecretsClient, err := deletedsecrets.NewDeletedSecretsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building DeletedSecrets client: %+v", err)
	}
	configureFunc(deletedSecretsClient.Client)

	deletedStorageClient, err := deletedstorage.NewDeletedStorageClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building DeletedStorage client: %+v", err)
	}
	configureFunc(deletedStorageClient.Client)

	fullBackupClient, err := fullbackup.NewFullBackupClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building FullBackup client: %+v", err)
	}
	configureFunc(fullBackupClient.Client)

	fullRestoreClient, err := fullrestore.NewFullRestoreClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building FullRestore client: %+v", err)
	}
	configureFunc(fullRestoreClient.Client)

	hSMSecurityDomainClient, err := hsmsecuritydomain.NewHSMSecurityDomainClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building HSMSecurityDomain client: %+v", err)
	}
	configureFunc(hSMSecurityDomainClient.Client)

	keysClient, err := keys.NewKeysClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Keys client: %+v", err)
	}
	configureFunc(keysClient.Client)

	rNGClient, err := rng.NewRNGClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building RNG client: %+v", err)
	}
	configureFunc(rNGClient.Client)

	roleAssignmentsClient, err := roleassignments.NewRoleAssignmentsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building RoleAssignments client: %+v", err)
	}
	configureFunc(roleAssignmentsClient.Client)

	roleDefinitionsClient, err := roledefinitions.NewRoleDefinitionsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building RoleDefinitions client: %+v", err)
	}
	configureFunc(roleDefinitionsClient.Client)

	secretsClient, err := secrets.NewSecretsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Secrets client: %+v", err)
	}
	configureFunc(secretsClient.Client)

	settingsClient, err := settings.NewSettingsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Settings client: %+v", err)
	}
	configureFunc(settingsClient.Client)

	storageClient, err := storage.NewStorageClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Storage client: %+v", err)
	}
	configureFunc(storageClient.Client)

	return &Client{
		Certificates:        certificatesClient,
		DeletedCertificates: deletedCertificatesClient,
		DeletedKeys:         deletedKeysClient,
		DeletedSecrets:      deletedSecretsClient,
		DeletedStorage:      deletedStorageClient,
		FullBackup:          fullBackupClient,
		FullRestore:         fullRestoreClient,
		HSMSecurityDomain:   hSMSecurityDomainClient,
		Keys:                keysClient,
		RNG:                 rNGClient,
		RoleAssignments:     roleAssignmentsClient,
		RoleDefinitions:     roleDefinitionsClient,
		Secrets:             secretsClient,
		Settings:            settingsClient,
		Storage:             storageClient,
	}, nil
}
