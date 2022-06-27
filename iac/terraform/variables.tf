variable "location" {
  description = "location of the resource"
  type        = string
  default     = "UK South"
}

variable "resources_prefix" {
  description = "prefix of resources that will be created from this automation"
  type        = string
  default     = null
}

