change "resource-enhancement" {
    body = "`azurerm_mongo_cluster` - the `storage_type` property is now optional and computed, so the value Azure chooses is kept when it is not set"
}

change "resource-enhancement" {
    body = "`azurerm_mongo_cluster` - the `storage_size_in_gb` property is now optional and computed, preserves the source cluster's inherited storage size for geo replicas and point-in-time restores when omitted"
}