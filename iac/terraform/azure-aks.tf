resource "azurerm_kubernetes_cluster" "aks" {
  name                = local.aks_cluster_name
  location            = azurerm_resource_group.aks.location
  dns_prefix          = local.aks_dns_prefix
  resource_group_name = azurerm_resource_group.aks.name


  default_node_pool {
    name                = "agentpool"
    enable_auto_scaling = true
    min_count           = 1
    max_count           = 3
    vm_size             = "Standard_DS2_v2"
    zones               = local.aks_zones
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "kubenet"
  }

  oms_agent {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.aks.id
  }

}

resource "azurerm_role_assignment" "aks_to_acr_role" {
  scope                            = azurerm_container_registry.container_registry.id
  role_definition_name             = "AcrPull"
  principal_id                     = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  skip_service_principal_aad_check = true
}