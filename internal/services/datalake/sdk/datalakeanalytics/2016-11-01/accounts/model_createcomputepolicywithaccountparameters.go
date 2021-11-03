package accounts

type CreateComputePolicyWithAccountParameters struct {
	Name       string                                `json:"name"`
	Properties CreateOrUpdateComputePolicyProperties `json:"properties"`
}
