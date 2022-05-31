package monitorsresource

type UserInfo struct {
	CompanyInfo  *CompanyInfo `json:"companyInfo,omitempty"`
	CompanyName  *string      `json:"companyName,omitempty"`
	EmailAddress *string      `json:"emailAddress,omitempty"`
	FirstName    *string      `json:"firstName,omitempty"`
	LastName     *string      `json:"lastName,omitempty"`
}
