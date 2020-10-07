variable "prefix" {
  description = "A prefix used for all resources in this example"
  default = "acisampleaks"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be provisioned"
  default = "West US"
}

variable "password" {
  description = "Service Principal password"
  default = "VT=uSgbTanZhyz@%nL9Hpd+Tfay_MRV#"
}
