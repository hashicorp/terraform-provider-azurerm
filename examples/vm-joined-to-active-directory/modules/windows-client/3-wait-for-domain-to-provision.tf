// NOTE: this is a hack.
// the AD Domain takes ~7m to provision, so we don't try and join an non-existant domain we sleep
// unfortunately we can't depend on the Domain Creation VM Extension since there's a reboot.
// We sleep for 12 minutes here to give Azure some breathing room.
resource "null_resource" "wait-for-domain-to-provision" {
  provisioner "local-exec" {
    command = "sleep 720"
  }

  depends_on = ["azurerm_virtual_machine.client"]
}
