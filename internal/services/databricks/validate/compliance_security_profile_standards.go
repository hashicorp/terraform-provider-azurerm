// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

type ComplianceStandard string

const (
	ComplianceStandardHIPAA  ComplianceStandard = "HIPAA"
	ComplianceStandardNONE   ComplianceStandard = "NONE"
	ComplianceStandardPCIDSS ComplianceStandard = "PCI_DSS"
)

func PossibleValuesForComplianceStandard() []string {
	return []string{
		string(ComplianceStandardHIPAA),
		string(ComplianceStandardPCIDSS),
	}
}
