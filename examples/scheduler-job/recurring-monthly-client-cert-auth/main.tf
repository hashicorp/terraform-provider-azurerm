resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_scheduler_job" "example" {
  name                = "${var.example}-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  action_web {
    url = "https://this.url.fails"

    authentication_certificate {
      pfx      = "${base64encode(file("your_cert.pfx"))}"
      password = "cert_password"
    }
  }

  error_action_web {
    url    = "https://this.url.fails"
    method = "put"
    body   = "The job failed"

    headers = {
      "Content-Type" = "text"
    }

    authentication_basic {
      username = "login"
      password = "apassword"
    }
  }

  recurrence {
    frequency = "monthly"
    count     = 1000

    monthly_occurrences = [
      {
        # first Sunday
        day        = "Sunday"
        occurrence = 1
      },
      {
        # third Sunday
        day        = "Sunday"
        occurrence = 3
      },
      {
        # last Sunday
        day        = "Sunday"
        occurrence = -1
      },
    ]
  }

  start_time = "2018-07-07T07:07:07-07:00"
}
