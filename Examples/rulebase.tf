resource "chkp_accessrulebaselist" "InitialRules" {
      name = "AzureScaleSetPackage"
      access = true
      threatprevention = false
      color = "pink"

          rulebase {
                enabled = true
                name = "Jump Host Rule"
                source = [ "${chkp_host.jumphost.name}","${chkp_host.ubuntu1.name}"]
                service = ["${data.chkp_servicetcp.ssh.name}"]
                action = "Accept"
                track {
                   type = "Log"
                }
                                }
          rulebase {
                name = "Internal Access"
                destination = ["${chkp_network.vnet_local.name}"]
                destinationnegate = false
                source = ["${chkp_network.vnet_local.name}"]
                enabled = true
                action = "Accept"
                track {
                      type = "Log"
                      }
                    }

          rulebase {
                name = "Internet Access"
                source = ["${chkp_network.vnet_local.name}","${chkp_securityzone.mysecurityzone.name}"]
                enabled = true
                action = "Accept"
                track {
                      type = "Log"
                      }
                    }

          rulebase {
                name = "Layer Test"
                source = ["${chkp_network.vnet_local.name}","${chkp_securityzone.mysecurityzone.name}"]
                enabled = true
                action = "Apply Layer"
                inlinelayer = "${chkp_accesslayerlist.LayerRules.name}"
                track {
                      type = "None"
                      perconnection = false
                      }
                    }
          rulebase {
                name = "Azure Health Checks"
                source = ["${chkp_host.azurelbhealthcheck.name}"]
                service = [ "${chkp_servicetcp.healthcheck.name}"]
                destination = ["${chkp_network.vnet_local.name}"]
                enabled = true
                action = "Accept"
                track {
                type = "Log"

                  }
                }

          rulebase {
                name = "AWS to Azure"
                destination = ["${chkp_network.vnet_local.name}"]
                source = ["${chkp_network.aws_VPC_1.name}"]
                service = ["${data.chkp_servicetcp.ssh.name}","${data.chkp_servicetcp.https.name}", "http"]
                enabled = true
                action = "Accept"
                track {
                      type = "Log"
                      }
                    }

          rulebase {
                name = "Azure to AWS"
                destination = ["${chkp_network.aws_VPC_1.name}"]
                source = ["${chkp_network.vnet_local.name}","${chkp_network.aws_VPC_1.name}"]
                service = ["${data.chkp_servicetcp.ssh.name}"]
                enabled = true
                action = "Accept"
                track {
                      type = "Log"
                      }
                    }

          rulebase {
                  name = "Cleanup Rule"
                  enabled = true
                  action = "Drop"
                  track {
                        type = "Log"
                        }
                      }
                  }
