resource "azurerm_kubernetes_cluster" "aks" {
  name                = local.aks_cluster_name
  location            = azurerm_resource_group.aks.location
  dns_prefix          = local.aks_dns_prefix
  resource_group_name = azurerm_resource_group.aks.name


  default_node_pool {
    name       = "agentpool"
    node_count = "1"
    vm_size    = "Standard_DS2_v2"
    zones      = ["1", "2"]
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