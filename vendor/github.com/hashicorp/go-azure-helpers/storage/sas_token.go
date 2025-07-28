// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	connStringAccountKeyKey    = "AccountKey"
	connStringAccountNameKey   = "AccountName"
	blobContainerSignedVersion = "2018-11-09"
)

// ComputeAccountSASToken computes the SAS Token for a Storage Account based on the
// access key & given permissions
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/create-account-sas
func ComputeAccountSASToken(accountName string,
	accountKey string,
	permissions string,
	services string,
	resourceTypes string,
	start string,
	expiry string,
	signedProtocol string,
	signedIp string, // nolint: unparam
	signedVersion string, // nolint: unparam
	signedEncryptionScope string, // nolint: unparam
) (string, error) {
	// UTF-8 by default...
	stringToSign := accountName + "\n"
	stringToSign += permissions + "\n"
	stringToSign += services + "\n"
	stringToSign += resourceTypes + "\n"
	stringToSign += start + "\n"
	stringToSign += expiry + "\n"
	stringToSign += signedIp + "\n"
	stringToSign += signedProtocol + "\n"
	stringToSign += signedVersion + "\n"

	if signedVersion >= "2020-12-06" {
		stringToSign += signedEncryptionScope + "\n"
	}

	binaryKey, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return "", err
	}
	hasher := hmac.New(sha256.New, binaryKey)
	hasher.Write([]byte(stringToSign))
	signature := hasher.Sum(nil)

	// Trial and error to determine which fields the Azure portal
	// URL encodes for a query string and which it does not.
	sasToken := "?sv=" + url.QueryEscape(signedVersion)
	sasToken += "&ss=" + url.QueryEscape(services)
	sasToken += "&srt=" + url.QueryEscape(resourceTypes)
	sasToken += "&sp=" + url.QueryEscape(permissions)
	sasToken += "&se=" + (expiry)
	sasToken += "&st=" + (start)
	sasToken += "&spr=" + (signedProtocol)

	// this is consistent with how the Azure portal builds these.
	if len(signedIp) > 0 {
		sasToken += "&sip=" + signedIp
	}

	sasToken += "&sig=" + url.QueryEscape(base64.StdEncoding.EncodeToString(signature))

	return sasToken, nil
}

// ComputeAccountSASConnectionString computes the composed SAS Connection String for a Storage Account based on the
// sas token
func ComputeAccountSASConnectionString(env *azure.Environment, accountName string, sasToken string) string {
	return fmt.Sprintf(
		"BlobEndpoint=https://%[1]s.blob.%[2]s/;"+
			"FileEndpoint=https://%[1]s.file.%[2]s/;"+
			"QueueEndpoint=https://%[1]s.queue.%[2]s/;"+
			"TableEndpoint=https://%[1]s.table.%[2]s/;"+
			"SharedAccessSignature=%[3]s", accountName, env.StorageEndpointSuffix, sasToken[1:]) // need to cut the first character '?' from the sas token
}

// ComputeAccountSASConnectionUrlForType computes the SAS Connection String for a Storage Account based on the
// sas token and the storage type
func ComputeAccountSASConnectionUrlForType(env *azure.Environment, accountName string, sasToken string, storageType string) (*string, error) {
	if !strings.EqualFold(storageType, "blob") && !strings.EqualFold(storageType, "file") && !strings.EqualFold(storageType, "queue") && !strings.EqualFold(storageType, "table") {
		return nil, fmt.Errorf("Unexpected storage type %s!", storageType)
	}

	url := fmt.Sprintf("https://%s.%s.%s%s", accountName, strings.ToLower(storageType), env.StorageEndpointSuffix, sasToken)
	return &url, nil
}

func ComputeContainerSASToken(signedPermissions string,
	signedStart string,
	signedExpiry string,
	accountName string,
	accountKey string,
	containerName string,
	signedIdentifier string,
	signedIp string,
	signedProtocol string,
	signedSnapshotTime string,
	cacheControl string,
	contentDisposition string,
	contentEncoding string,
	contentLanguage string,
	contentType string,
) (string, error) {
	canonicalizedResource := "/blob/" + accountName + "/" + containerName
	signedVersion := blobContainerSignedVersion
	signedResource := "c" // c for container

	// UTF-8 by default...
	stringToSign := signedPermissions + "\n"
	stringToSign += signedStart + "\n"
	stringToSign += signedExpiry + "\n"
	stringToSign += canonicalizedResource + "\n"
	stringToSign += signedIdentifier + "\n"
	stringToSign += signedIp + "\n"
	stringToSign += signedProtocol + "\n"
	stringToSign += signedVersion + "\n"
	stringToSign += signedResource + "\n"
	stringToSign += signedSnapshotTime + "\n"
	stringToSign += cacheControl + "\n"
	stringToSign += contentDisposition + "\n"
	stringToSign += contentEncoding + "\n"
	stringToSign += contentLanguage + "\n"
	stringToSign += contentType

	binaryKey, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return "", err
	}
	hasher := hmac.New(sha256.New, binaryKey)
	hasher.Write([]byte(stringToSign))
	signature := hasher.Sum(nil)

	sasToken := "?sv=" + signedVersion
	sasToken += "&sr=" + signedResource
	sasToken += "&st=" + signedStart
	sasToken += "&se=" + signedExpiry
	sasToken += "&sp=" + signedPermissions

	if len(signedIp) > 0 {
		sasToken += "&sip=" + signedIp
	}

	if len(signedProtocol) > 0 {
		sasToken += "&spr=" + signedProtocol
	}

	if len(signedIdentifier) > 0 {
		sasToken += "&si=" + signedIdentifier
	}

	if len(cacheControl) > 0 {
		sasToken += "&rscc=" + url.QueryEscape(cacheControl)
	}

	if len(contentDisposition) > 0 {
		sasToken += "&rscd=" + url.QueryEscape(contentDisposition)
	}

	if len(contentEncoding) > 0 {
		sasToken += "&rsce=" + url.QueryEscape(contentEncoding)
	}

	if len(contentLanguage) > 0 {
		sasToken += "&rscl=" + url.QueryEscape(contentLanguage)
	}

	if len(contentType) > 0 {
		sasToken += "&rsct=" + url.QueryEscape(contentType)
	}

	sasToken += "&sig=" + url.QueryEscape(base64.StdEncoding.EncodeToString(signature))

	return sasToken, nil
}

// ParseAccountSASConnectionString parses the Connection String for a Storage Account
func ParseAccountSASConnectionString(connString string) (map[string]string, error) {
	// This connection string was for a real storage account which has been deleted
	// so its safe to include here for reference to understand the format.
	// DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==;EndpointSuffix=core.windows.net
	validKeys := map[string]bool{
		"DefaultEndpointsProtocol": true, "BlobEndpoint": true,
		"AccountName": true, "AccountKey": true, "EndpointSuffix": true,
	}
	// The k-v pairs are separated with semi-colons
	tokens := strings.Split(connString, ";")

	kvp := make(map[string]string)

	for _, atoken := range tokens {
		// The individual k-v are separated by an equals sign.
		kv := strings.SplitN(atoken, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("[ERROR] token `%s` is an invalid key=pair (connection string %s)", atoken, connString)
		}

		key := kv[0]
		val := kv[1]

		if _, present := validKeys[key]; !present {
			return nil, fmt.Errorf("[ERROR] Unknown Key `%s` in connection string %s", key, connString)
		}
		kvp[key] = val
	}

	if _, present := kvp[connStringAccountKeyKey]; !present {
		return nil, fmt.Errorf("[ERROR] Storage Account Key not found in connection string: %s", connString)
	}

	return kvp, nil
}
