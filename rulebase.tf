resource "chkp_accessrulebaselist" "InitialRules" {
      name = "AzureScaleSetPackage"
      access = true
      threatprevention = false
      color = "pink"

rulebase {
      enabled = true
      name = "Jump Host Rule"
      source = [ "${chkp_host.jumphost.uid}","${chkp_host.ubuntu1.id}"]
      service = ["${data.chkp_servicetcp.ssh.id}"]
      action = "accept"
      track {
         type = "Log"
      }
                      }
rulebase {
      name = "Internal Access"
      destination = ["${chkp_network.vnet_local.id}"]
      destinationnegate = false
      source = ["${chkp_network.vnet_local.id}"]
      enabled = true
      action = "drop"
      track {
            type = "Log"
            }
          }

rulebase {
      name = "Internet Access"
      source = ["${chkp_network.vnet_local.id}"]
      enabled = true
      action = "accept"
      track {
            type = "Log"
            }
          }

rulebase {
      name = "Azure Health Checks"
      source = ["${chkp_host.azurelbhealthcheck.uid}"]
      service = [ "${chkp_servicetcp.healthcheck.id}"]
      destination = ["${chkp_network.vnet_local.id}"]
      enabled = true
      action = "accept"
      track {
      type = "Log"
        }
      }

rulebase {
      name = "AWS to Azure"
      destination = ["${chkp_network.vnet_local.id}"]
      source = ["${chkp_network.aws_VPC_1.uid}"]
      service = ["${data.chkp_servicetcp.ssh.id}","${data.chkp_servicetcp.https.id}"]
      enabled = true
      action = "drop"
      track {
            type = "Log"
            }
          }

rulebase {
      name = "Azure to AWS"
      destination = ["${chkp_network.aws_VPC_1.uid}"]
      source = ["${chkp_network.vnet_local.id}","${chkp_network.aws_VPC_1.uid}"]
      service = ["${data.chkp_servicetcp.ssh.id}"]
      enabled = true
      action = "accept"
      track {
            type = "Log"
            }
          }

rulebase {
        name = "Cleanup Rule"
        enabled = true
        action = "drop"
        track {
              type = "Log"
              }
            }
        }
