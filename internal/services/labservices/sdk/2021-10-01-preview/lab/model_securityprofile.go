package lab

type SecurityProfile struct {
	OpenAccess       *EnableState `json:"openAccess,omitempty"`
	RegistrationCode *string      `json:"registrationCode,omitempty"`
}
