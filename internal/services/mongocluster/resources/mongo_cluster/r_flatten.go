// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"net/url"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
)

func flattenMongoClusterPreviewFeatures(input *[]mongoclusters.PreviewFeature) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func flattenMongoClusterConnectionStrings(input *[]mongoclusters.ConnectionString, userName, userPassword string) []MongoClusterConnectionString {
	results := make([]MongoClusterConnectionString, 0)
	if input == nil {
		return results
	}
	for _, cs := range *input {
		value := pointer.From(cs.ConnectionString)
		// Password can be empty if it isn't available in the state file (e.g. during import).
		// In this case, we simply leave the placeholder unchanged.
		if userPassword != "" {
			value = regexp.MustCompile(`<user>:<password>`).ReplaceAllLiteralString(value, url.UserPassword(userName, userPassword).String())
		}

		results = append(results, MongoClusterConnectionString{
			Name:        pointer.From(cs.Name),
			Description: pointer.From(cs.Description),
			Value:       value,
		})
	}

	return results
}

func flattenMongoClusterCustomerManagedKey(input *mongoclusters.EncryptionProperties) []CustomerManagedKey {
	if input == nil || input.CustomerManagedKeyEncryption == nil {
		return []CustomerManagedKey{}
	}

	cmkEncryption := input.CustomerManagedKeyEncryption
	uaiResourceId := ""
	if cmkEncryption.KeyEncryptionKeyIdentity != nil {
		uaiResourceId = pointer.From(cmkEncryption.KeyEncryptionKeyIdentity.UserAssignedIdentityResourceId)
	}

	return []CustomerManagedKey{
		{
			KeyVaultKeyId:          pointer.From(cmkEncryption.KeyEncryptionKeyURL),
			UserAssignedIdentityId: uaiResourceId,
		},
	}
}

func flattenMongoClusterRestore(input *mongoclusters.MongoClusterRestoreParameters) []Restore {
	results := make([]Restore, 0)
	if input == nil {
		return results
	}

	results = append(results, Restore{
		PointInTimeUtc: pointer.From(input.PointInTimeUTC),
		SourceId:       pointer.From(input.SourceResourceId),
	})

	return results
}

func flattenMongoClusterAuthConfig(input *mongoclusters.AuthConfigProperties) []string {
	results := make([]string, 0)
	if input == nil || input.AllowedModes == nil {
		return results
	}

	for _, v := range *input.AllowedModes {
		results = append(results, string(v))
	}

	return results
}
