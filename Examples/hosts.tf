resource "chkp_host" "jumphost" {
  ipv4address = "192.168.198.198"
  name = "jumphost"
  color = "blue"
 # natsettings {
 #   autorule = true
 #   hidebehind = "gateway"
 #   method = "hide"
 #   installon = "Test_GW"
 # }
    }

resource "chkp_host" "ubuntu1" {
      ipv4address = "172.16.32.10"
      name = "ubuntu1"
      color = "blue"
      }

resource "chkp_host" "azurelbhealthcheck" {
      ipv4address = "168.63.129.16"
      name = "azure_lb_health_check"
      color = "light green"
      }
