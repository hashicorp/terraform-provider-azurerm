package monitorsresource

type CompanyInfo struct {
	Business        *string `json:"business,omitempty"`
	Country         *string `json:"country,omitempty"`
	Domain          *string `json:"domain,omitempty"`
	EmployeesNumber *string `json:"employeesNumber,omitempty"`
	State           *string `json:"state,omitempty"`
}
