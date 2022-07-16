resource "azurerm_dns_zone" "dns_zone" {
  name                = local.dns_zone
  resource_group_name = azurerm_resource_group.aks.name
}