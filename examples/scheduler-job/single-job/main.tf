resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_scheduler_job" "example" {
  name                = "${var.prefix}-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  # re-enable it each run
  state = "enabled"

  action_web {
    # defaults to get
    url = "http://this.url.fails"
  }

  # default start time is now
}
