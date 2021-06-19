package dataprotection

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackupVault -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/vault1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackupPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
