variable "description" {
  type = string
  description = "Description of the scheduled query rule."
}

variable "enabled" {
  type = bool
  description = "If the scheduled query rule is enabled."
}

variable "query" {
  type = string
  description = "Log search query. Required for action_type - AlertingAction."
}

variable "authorized_resources" {
  type = list(string)
  description = "List of resources referred into query."
}

variable "data_source_id" {
  type = string
  description = "The resource uri over which log search query is to be run."
}

variable "schedule" {
  type = object({
    frequency_in_minutes = number
    time_window_in_minutes = number
  })
  description = "frequency_in_minutes - frequency (in minutes) at which rule condition should be evaluated. time_window_in_minutes - Time window for which data needs to be fetched for query (should be greater than or equal to frequencyInMinutes)."
}

variable "action" {
  type = string
  default = "Action"
  description = "Must be one of 'Action', 'AlertingAction', or 'LogToMetricAction'."
}
