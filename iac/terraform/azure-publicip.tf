resource "azurerm_public_ip" "aks" {
  name                = local.aks_public_ip
  resource_group_name = azurerm_resource_group.aks.node_resource_group
  location            = azurerm_resource_group.aks.location
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = local.aks_zones
  ip_version          = "IPv4"
}