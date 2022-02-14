package kubernetes

type CredentialResults struct {
	HybridConnectionConfig *HybridConnectionConfig `json:"hybridConnectionConfig,omitempty"`
	Kubeconfigs            *[]CredentialResult     `json:"kubeconfigs,omitempty"`
}
