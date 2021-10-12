package lab

type RosterProfile struct {
	ActiveDirectoryGroupId *string `json:"activeDirectoryGroupId,omitempty"`
	LmsInstance            *string `json:"lmsInstance,omitempty"`
	LtiClientId            *string `json:"ltiClientId,omitempty"`
	LtiContextId           *string `json:"ltiContextId,omitempty"`
	LtiRosterEndpoint      *string `json:"ltiRosterEndpoint,omitempty"`
}
