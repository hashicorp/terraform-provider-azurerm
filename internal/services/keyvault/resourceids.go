// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

// KeyVault Access Policies are Terraform specific, but can be either an Object ID or an Application ID
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AccessPolicyApplication -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/object1/applicationId/application1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AccessPolicyObject -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/object1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Key -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/keys/key1/versions/version1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=KeyVersionless -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/keys/key1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Secret -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/secrets/secret1/versions/version1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SecretVersionless -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/secrets/secret1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Certificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/certificates/cert1/versions/version1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CertificateVersionless -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/certificates/cert1
