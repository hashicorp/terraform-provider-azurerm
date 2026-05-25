// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

type ComplianceStandard string

const (
	ComplianceStandardHIPAA              ComplianceStandard = "HIPAA"
	ComplianceStandardNONE               ComplianceStandard = "NONE"
	ComplianceStandardPCIDSS             ComplianceStandard = "PCI_DSS"
	ComplianceStandardFEDRAMPMODERATE    ComplianceStandard = "FEDRAMP_MODERATE"
	ComplianceStandardIRAPPROTECTED      ComplianceStandard = "IRAP_PROTECTED"
	ComplianceStandardFEDRAMPHIGH        ComplianceStandard = "FEDRAMP_HIGH"
	ComplianceStandardFEDRAMPIL5         ComplianceStandard = "FEDRAMP_IL5"
	ComplianceStandardITAREAR            ComplianceStandard = "ITAR_EAR"
	ComplianceStandardCYBERESSENTIALPLUS ComplianceStandard = "CYBER_ESSENTIAL_PLUS"
	ComplianceStandardCANADAPROTECTEDB   ComplianceStandard = "CANADA_PROTECTED_B"
	ComplianceStandardISMAP              ComplianceStandard = "ISMAP"
	ComplianceStandardHITRUST            ComplianceStandard = "HITRUST"
	ComplianceStandardKFSI               ComplianceStandard = "K_FSI"
	ComplianceStandardGERMANYC5          ComplianceStandard = "GERMANY_C5"
	ComplianceStandardGERMANYTISAX       ComplianceStandard = "GERMANY_TISAX"
)

func PossibleValuesForComplianceStandard() []string {
	return []string{
		string(ComplianceStandardHIPAA),
		string(ComplianceStandardPCIDSS),
		string(ComplianceStandardFEDRAMPMODERATE),
		string(ComplianceStandardIRAPPROTECTED),
		string(ComplianceStandardFEDRAMPHIGH),
		string(ComplianceStandardFEDRAMPIL5),
		string(ComplianceStandardITAREAR),
		string(ComplianceStandardCYBERESSENTIALPLUS),
		string(ComplianceStandardCANADAPROTECTEDB),
		string(ComplianceStandardISMAP),
		string(ComplianceStandardHITRUST),
		string(ComplianceStandardKFSI),
		string(ComplianceStandardGERMANYC5),
		string(ComplianceStandardGERMANYTISAX),
	}
}
