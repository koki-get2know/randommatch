locals {
  resources_prefix    = var.resources_prefix != null ? var.resources_prefix : "${local._default.name_prefix}"
  location            = var.location
  resource_group_name = "${local.resources_prefix}rg"
}