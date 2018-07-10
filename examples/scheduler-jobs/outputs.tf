output "job_collection-id" {
  value = "${azurerm_scheduler_job_collection.example.id}"
}

output "job-web-once-url" {
  value = "${azurerm_scheduler_job.web-once-now.action_web.0.url}"
}

