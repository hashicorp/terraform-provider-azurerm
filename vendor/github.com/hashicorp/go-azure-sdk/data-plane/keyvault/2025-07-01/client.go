package v2025_07_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/backuprestores"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/certificates"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/deletedstorage"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/keys"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/rbacs"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/secrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/securitydomains"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/settings"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/storage"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

type Client struct {
	Backuprestores  *backuprestores.BackuprestoresClient
	Certificates    *certificates.CertificatesClient
	DeletedStorage  *deletedstorage.DeletedStorageClient
	Keys            *keys.KeysClient
	Rbacs           *rbacs.RbacsClient
	Secrets         *secrets.SecretsClient
	Securitydomains *securitydomains.SecuritydomainsClient
	Settings        *settings.SettingsClient
	Storage         *storage.StorageClient
}

func NewClient(configureFunc func(c *dataplane.Client)) (*Client, error) {
	backuprestoresClient, err := backuprestores.NewBackuprestoresClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Backuprestores client: %+v", err)
	}
	configureFunc(backuprestoresClient.Client)

	certificatesClient, err := certificates.NewCertificatesClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Certificates client: %+v", err)
	}
	configureFunc(certificatesClient.Client)

	deletedStorageClient, err := deletedstorage.NewDeletedStorageClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building DeletedStorage client: %+v", err)
	}
	configureFunc(deletedStorageClient.Client)

	keysClient, err := keys.NewKeysClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Keys client: %+v", err)
	}
	configureFunc(keysClient.Client)

	rbacsClient, err := rbacs.NewRbacsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Rbacs client: %+v", err)
	}
	configureFunc(rbacsClient.Client)

	secretsClient, err := secrets.NewSecretsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Secrets client: %+v", err)
	}
	configureFunc(secretsClient.Client)

	securitydomainsClient, err := securitydomains.NewSecuritydomainsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Securitydomains client: %+v", err)
	}
	configureFunc(securitydomainsClient.Client)

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
		Backuprestores:  backuprestoresClient,
		Certificates:    certificatesClient,
		DeletedStorage:  deletedStorageClient,
		Keys:            keysClient,
		Rbacs:           rbacsClient,
		Secrets:         secretsClient,
		Securitydomains: securitydomainsClient,
		Settings:        settingsClient,
		Storage:         storageClient,
	}, nil
}
