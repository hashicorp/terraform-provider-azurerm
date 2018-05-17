output "job_collection-id" {
  value = "${azurerm_scheduler_job_collection.jc.id}"
}

output "job-web-once-url" {
  value = "${azurerm_scheduler_job.web-once-now.action_web.0.url}"
}

output "job-web-once-count" {
  value = "${azurerm_scheduler_job.web-once-now.execution_count}"
}

output "job-web-recurring_weekly-authCert_thumb" {
  value = "${azurerm_scheduler_job.web-recurring_weekly-auth_cert.action_web.0.authentication_certificate.0.thumbprint}"
}

output "job-web-recurring_weekly-authCert_expire" {
  value = "${azurerm_scheduler_job.web-recurring_weekly-auth_cert.action_web.0.authentication_certificate.0.expiration}"
}

output "job-web-recurring_weekly-authCert_san" {
  value = "${azurerm_scheduler_job.web-recurring_weekly-auth_cert.action_web.0.authentication_certificate.0.subject_name}"
}
