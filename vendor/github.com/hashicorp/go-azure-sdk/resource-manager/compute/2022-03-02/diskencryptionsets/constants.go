package diskencryptionsets

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionSetType string

const (
	DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey      DiskEncryptionSetType = "ConfidentialVmEncryptedWithCustomerKey"
	DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey             DiskEncryptionSetType = "EncryptionAtRestWithCustomerKey"
	DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys DiskEncryptionSetType = "EncryptionAtRestWithPlatformAndCustomerKeys"
)

func PossibleValuesForDiskEncryptionSetType() []string {
	return []string{
		string(DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey),
		string(DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
		string(DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys),
	}
}

func parseDiskEncryptionSetType(input string) (*DiskEncryptionSetType, error) {
	vals := map[string]DiskEncryptionSetType{
		"confidentialvmencryptedwithcustomerkey":      DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey,
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
