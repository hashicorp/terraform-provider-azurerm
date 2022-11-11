package diskencryptionsets

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionSetType string

const (
	DiskEncryptionSetTypeConfidentialVmEncryptedWithCustomerKey      DiskEncryptionSetType = "ConfidentialVmEncryptedWithCustomerKey"
	DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey             DiskEncryptionSetType = "EncryptionAtRestWithCustomerKey"
	DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys DiskEncryptionSetType = "EncryptionAtRestWithPlatformAndCustomerKeys"
)

func PossibleValuesForDiskEncryptionSetType() []string {
	return []string{
		string(DiskEncryptionSetTypeConfidentialVmEncryptedWithCustomerKey),
		string(DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
		string(DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys),
	}
}

func parseDiskEncryptionSetType(input string) (*DiskEncryptionSetType, error) {
	vals := map[string]DiskEncryptionSetType{
		"confidentialvmencryptedwithcustomerkey":      DiskEncryptionSetTypeConfidentialVmEncryptedWithCustomerKey,
		"encryptionatrestwithcustomerkey":             DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey,
		"encryptionatrestwithplatformandcustomerkeys": DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskEncryptionSetType(input)
	return &out, nil
}
