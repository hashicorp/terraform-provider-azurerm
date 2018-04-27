
resource "azurerm_resource_group" "rg" {
    name     = "${var.resource_group_name}"
    location = "${var.resource_group_location}"
}

resource "azurerm_scheduler_job_collection" "jc" {
    name                = "tfex-scheduler-job-collection"
    location            = "${azurerm_resource_group.rg.location}"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    sku                 = "standard"
    state               = "enabled"

    quota {
        max_job_count            = 15
        max_recurrence_interval  = 1
        max_recurrence_frequency = "minute"
    }
}

resource "azurerm_scheduler_job" "web-once-now" {
    name                = "tfex-web-once-now"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    job_collection_name = "${azurerm_scheduler_job_collection.jc.name}"

    state = "enabled"

    action_web {
        url    = "http://this.url.fails"
    }

    retry {
        interval = "00:05:00" //5 min
        count    =  4
    }

    //times in the past start imediatly and run once,
    start_time = "1987-07-07T07:07:07-07:00"
}

resource "azurerm_scheduler_job" "web-recurring" {
    name                = "tfex-web-recurring"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    job_collection_name = "${azurerm_scheduler_job_collection.jc.name}"

    action_web {
        url     = "http://this.url.fails"
        method  = "put"
        body    = "this is some text"
        headers = {
            Content-Type = "text"
        }
    }

    retry {
        interval = "00:05:00" //5 min
        count    =  4
    }

    recurrence {
        frequency  = "minute"
        interval   = 5
        count      = 10 //recurring counts start in past
    }

    start_time = "2019-07-07T07:07:07-07:00"
}

resource "azurerm_scheduler_job" "web-recurring_daily-auth_basic" {
    name                = "tfex-web-recurring_daily-auth_basic"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    job_collection_name = "${azurerm_scheduler_job_collection.jc.name}"

    state = "enabled"

    action_web {
        url    = "https://this.url.fails"
        method = "get"

        authentication_basic {
            username = "login"
            password = "apassword"
        }
    }

    recurrence {
        frequency = "day"
        count     = 1000
        hours     = [0,12]
        minutes   = [0,15,30,45]
    }

    start_time = "2019-07-07T07:07:07-07:00"
}


resource "azurerm_scheduler_job" "web-recurring_weekly-auth_cert" {
    name                = "tfex-web-recurring_weekly-auth_cert"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    job_collection_name = "${azurerm_scheduler_job_collection.jc.name}"

    state = "enabled"

    action_web {
        url    = "https://this.url.fails"
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


resource "azurerm_scheduler_job" "web-recurring_monthly-error_action" {
    name                = "tfex-web-recurring_monthly-auth_ad"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    job_collection_name = "${azurerm_scheduler_job_collection.jc.name}"

    state = "enabled"

    action_web {
        url    = "http://this.url.fails"
        method = "get"
    }

    error_action_web {
        url     = "https://this.url.fails"
        method  = "put"
        body    = "The job failed"

        headers = {
            Content-Type = "text"
        }

        authentication_basic {
            username = "login"
            password = "apassword"
        }
    }

    recurrence {
        frequency    = "month"
        count        = 1000
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
            }
        ]
    }

    start_time = "2019-07-07T07:07:07-07:00"
}