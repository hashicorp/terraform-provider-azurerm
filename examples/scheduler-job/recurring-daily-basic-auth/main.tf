resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_scheduler_job" "example" {
  name                = "${var.prefix}-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  action_web {
    url    = "https://this.url.fails"
    method = "put"
    body   = "this is some text"

    headers = {
      Content-Type = "text"
    }

    authentication_basic {
      username = "login"
      password = "apassword"
    }
  }

  retry {
    # retry every 5 min a maximum of 10 times
    interval = "00:05:00"
    count    = 10
  }

  recurrence {
    frequency = "day"
    count     = 1000

    # run 4 times an hour every 12 hours
    hours   = [0, 12]
    minutes = [0, 15, 30, 45]
  }

  start_time = "2018-07-07T07:07:07-07:00"
}
