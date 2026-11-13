// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestResourceGuardVaultCriticalOperationExclusionList(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "Microsoft.RecoveryServices/vaults/backupconfig/delete",
			expected: true,
		},
		{
			input:    "Microsoft.RecoveryServices/vaults/backupResourceGuardProxies/write",
			expected: true,
		},
		{
			input:    "Microsoft.RecoveryServices/vaults/write#reduceImmutabilityState",
			expected: true,
		},
		{
			input:    "Microsoft.RecoveryServices/vaults/write#modifyEncryptionSettings",
			expected: true,
		},
		{
			input:    "Microsoft.DataProtection/backupVaults/backupInstances/stopProtection/action",
			expected: true,
		},
		{
			input:    "Microsoft.DataProtection/backupVaults/write#reduceImmutabilityState",
			expected: true,
		},
		{
			input:    "Microsoft.DataProtection/backupVaults/write#modifyEncryptionSettings",
			expected: true,
		},
		{
			input:    "Microsoft.RecoveryServices/vaults/backupconfig/write",
			expected: false,
		},
		{
			input:    "Microsoft.DataProtection/backupVaults/backupResourceGuardProxies/delete",
			expected: false,
		},
		{
			input:    "Microsoft.DataProtection/backupVaults/write#reduceSoftDeleteSecurity",
			expected: false,
		},
		{
			input:    "Microsoft.RecoveryServices/vaults/backupPolicies/delete",
			expected: false,
		},
		{
			input:    "Microsoft.DataProtection/backupVaults/backupPolicies/write",
			expected: false,
		},
		{
			input:    "Microsoft.Compute/virtualMachines/delete",
			expected: false,
		},
	}

	validator := ResourceGuardVaultCriticalOperationExclusionList()

	for _, testCase := range testCases {
		_, errors := validator(testCase.input, "vault_critical_operation_exclusion_list")
		result := len(errors) == 0
		if result != testCase.expected {
			t.Fatalf("expected the result to be %t but got %t for %q", testCase.expected, result, testCase.input)
		}
	}
}
