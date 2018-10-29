resource "azurerm_resource_group" "example" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

resource "azurerm_scheduler_job_collection" "example" {
  name                = "tfex-scheduler-job-collection"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "standard"
  state               = "enabled"

  quota {
    max_job_count            = 10
    max_recurrence_interval  = 10
    max_recurrence_frequency = "minute"
  }
}

resource "azurerm_scheduler_job" "web-once-now" {
  name                = "tfex-web-once-now"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  state = "enabled"

  action_web {
    url    = "http://example.com"
    method = "get"
  }

  retry {
    interval = "00:05:00" //5 min
    count    = 4
  }

  //times in the past start immediatly and run once,
  start_time = "1987-07-07T07:07:07-07:00"
}

resource "azurerm_scheduler_job" "web-recurring" {
  name                = "tfex-web-recurring"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  action_web {
    url    = "http://example.com"
    method = "put"
    body   = "this is some text"

    headers = {
      Content-Type = "text"
    }
  }

  retry {
    interval = "00:05:00" //5 min
    count    = 4
  }

  recurrence {
    frequency = "minute"
    interval  = 10
    count     = 10       //recurring counts start in past
  }

  start_time = "2019-07-07T07:07:07-07:00"
}

resource "azurerm_scheduler_job" "web-recurring_daily-auth_basic" {
  name                = "tfex-web-recurring_daily-auth_basic"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  state = "enabled"

  action_web {
    url    = "https://example.com"
    method = "get"

    authentication_basic {
      username = "login"
      password = "apassword"
    }
  }

  recurrence {
    frequency = "day"
    count     = 1000
    hours     = [0, 12]
    minutes   = [0, 15, 30, 45]
  }

  start_time = "2019-07-07T07:07:07-07:00"
}

resource "azurerm_scheduler_job" "web-recurring_weekly-auth_cert" {
  name                = "tfex-web-recurring_weekly-auth_cert"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  state = "enabled"

  action_web {
    url    = "https://example.com"
    method = "get"

    authentication_certificate {
      pfx      = "${base64encode(file("../../azurerm/testdata/application_gateway_test.pfx"))}"
      password = "terraform"
    }
  }

  recurrence {
    frequency = "week"
    count     = 1000
    week_days = ["Sunday", "Saturday"]
  }

  start_time = "2019-07-07T07:07:07-07:00"
}

resource "azurerm_storage_account" "example" {
  name                     = "tfexstorageaccount"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "example" {
  name                 = "tfex-schedulerjob-storagequeue"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  storage_account_name = "${azurerm_storage_account.example.name}"
}

resource "azurerm_scheduler_job" "storage-queue" {
  name                = "tfex-schedulerjob-storage_queue"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  action_storage_queue = {
    storage_account_name = "${azurerm_storage_account.example.name}"
    storage_queue_name   = "${azurerm_storage_queue.example.name}"
    sas_token            = "${azurerm_storage_account.example.primary_access_key}"
    message              = "storage message"
  }

  recurrence {
    frequency = "week"
    count     = 1000
    week_days = ["Sunday"]
  }

  start_time = "2019-07-07T07:07:07-07:00"
}

resource "azurerm_scheduler_job" "web-recurring_monthly-error_action" {
  name                = "tfex-web-recurring_monthly-auth_ad"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  state = "enabled"

  action_web {
    url    = "http://example.com"
    method = "get"
  }

  error_action_web {
    url    = "https://example.com"
    method = "put"
    body   = "The job failed"

    headers = {
      Content-Type = "text"
    }

    authentication_basic {
      username = "login"
      password = "apassword"
    }
  }

  recurrence {
    frequency = "month"
    count     = 1000

    monthly_occurrences = [
      {
        day        = "Sunday"
        occurrence = 1
      },
      {
        day        = "Sunday"
        occurrence = 3
      },
      {
        day        = "Sunday"
        occurrence = -1
      },
    ]
  }

  start_time = "2019-07-07T07:07:07-07:00"
}
