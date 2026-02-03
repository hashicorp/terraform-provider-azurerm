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
		string(ComplianceStandardNONE),
		string(ComplianceStandardPCIDSS),
	}
}
