resource "chkp_network" "vnet_local" {
      subnet4 = "10.99.0.0"
      masklength4 = "16"
      name = "Azure-Local-vNet"
      color = "blue"
      natsettings {
        autorule = true
        hidebehind = "gateway"
        method = "hide"
        installon = "All"
      }
        }

resource "chkp_network" "vnet_remote" {
      subnet4 = "10.199.0.0"
      masklength4 = "16"
      name = "Azure-Remote-vNet"
      color = "blue"
      }

resource "chkp_network" "aws_VPC_1" {
      subnet4 = "10.198.0.0"
      masklength4 = "16"
      name = "aws_VPC_1"
      color = "orange"
      }
