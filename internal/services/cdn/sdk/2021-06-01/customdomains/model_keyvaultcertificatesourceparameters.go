package customdomains

type KeyVaultCertificateSourceParameters struct {
	DeleteRule        DeleteRule `json:"deleteRule"`
	ResourceGroupName string     `json:"resourceGroupName"`
	SecretName        string     `json:"secretName"`
	SecretVersion     *string    `json:"secretVersion,omitempty"`
	SubscriptionId    string     `json:"subscriptionId"`
	TypeName          TypeName   `json:"typeName"`
	UpdateRule        UpdateRule `json:"updateRule"`
	VaultName         string     `json:"vaultName"`
}
