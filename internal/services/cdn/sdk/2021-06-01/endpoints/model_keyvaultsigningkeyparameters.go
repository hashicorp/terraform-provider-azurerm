package endpoints

type KeyVaultSigningKeyParameters struct {
	ResourceGroupName string   `json:"resourceGroupName"`
	SecretName        string   `json:"secretName"`
	SecretVersion     string   `json:"secretVersion"`
	SubscriptionId    string   `json:"subscriptionId"`
	TypeName          TypeName `json:"typeName"`
	VaultName         string   `json:"vaultName"`
}
