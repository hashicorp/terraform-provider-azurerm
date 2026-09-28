change "resource-enhancement" {
  body = "`azurerm_netapp_volume_group_oracle` - allow `volume.storage_quota_in_gb` minimum of 50 (previously 100)"
}

change "resource-enhancement" {
  body = "`azurerm_netapp_volume_group_sap_hana` - allow `volume.storage_quota_in_gb` minimum of 50 (previously 100)"
}
