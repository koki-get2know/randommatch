resource "azurerm_dns_zone" "dns_zone" {
  name                = local.dns_zone
  resource_group_name = azurerm_resource_group.aks.name
}

/*
resource "azurerm_dns_a_record" "dns_record" {
  name                = "koki"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  target_resource_id  = azurerm_public_ip.example.id
}*/