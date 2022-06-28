locals {
  resources_prefix             = var.resources_prefix != null ? var.resources_prefix : "${local._default.name_prefix}"
  location                     = var.location
  resource_group_name          = "${local.resources_prefix}aksrg"
  container_registry_name      = "${local.resources_prefix}akscr"
  log_analytics_workspace_name = "${local.resources_prefix}aksmonitor"
  log_analytics_workspace_sku  = "PerGB2018"
  aks_cluster_name             = "${local.resources_prefix}aksk8s"
  aks_dns_prefix               = "ps-${local.aks_cluster_name}"
}