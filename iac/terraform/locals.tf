locals {
  resources_prefix             = var.resources_prefix != null ? var.resources_prefix : "${local._default.name_prefix}"
  location                     = var.location
  resource_group_name          = "${local.resources_prefix}aksrg"
  container_registry_name      = "${local.resources_prefix}3akscr"
  log_analytics_workspace_name = "${local.resources_prefix}aksmonitor"
  log_analytics_workspace_sku  = "PerGB2018"
  aks_cluster_name             = "${local.resources_prefix}aksk8s"
  aks_dns_prefix               = "ps-${local.aks_cluster_name}"
  dns_zone                     = var.dns_zone
  aks_public_ip                = "${local.resources_prefix}aksstaticpublicip"
  aks_zones                    = ["1", "2", "3"]
  zone_prefix_name             = "koki"
  tenant_owner_object_id       = "122fe06e-afd1-4a13-91a9-975f904da0e9"
}
