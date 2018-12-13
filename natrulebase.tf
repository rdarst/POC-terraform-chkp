resource "chkp_accessnatlist" "InitialNATRules" {
      package = "${chkp_accessrulebaselist.InitialRules.uid}"

          natlist {
                enabled = true
                originalsource = "${chkp_host.jumphost.name}"
                originaldestination = "${chkp_host.ubuntu1.name}"
                originalservice = "http"
                translatedservice = "${data.chkp_servicetcp.ssh.name}"
                translatedsource = "${chkp_host.jumphost.name}"
                translateddestination = "${chkp_host.ubuntu1.name}"
                comments = "nat rule 1"
                enabled = true
                method = "hide"
                                }
          natlist {
                enabled = true
                originalsource = "${chkp_host.jumphost.name}"
                originaldestination = "${chkp_host.ubuntu1.name}"
                originalservice = "https"
                translatedservice = "smtp"
                translatedsource = "${chkp_host.jumphost.name}"
                translateddestination = "${chkp_host.ubuntu1.name}"
                comments = "nat rule 2"
                enabled = true
                method = "hide"
                                }

            natlist {
                  enabled = true
                  originalsource = "${chkp_host.jumphost.name}"
                  originaldestination = "${chkp_host.ubuntu1.name}"
                  originalservice = "ftp"
                  translatedservice = "http"
                  translatedsource = "${chkp_host.jumphost.name}"
                  translateddestination = "${chkp_host.ubuntu1.name}"
                  comments = "nat rule 3"
                  enabled = true
                  method = "static"
                                  }

            natlist {
                  enabled = true
                  originalsource = "${chkp_host.jumphost.name}"
                  originaldestination = "${chkp_host.ubuntu1.name}"
                  originalservice = "smtp"
                  translatedservice = "${data.chkp_servicetcp.ssh.name}"
                  translatedsource = "${chkp_host.jumphost.name}"
                  translateddestination = "${chkp_host.ubuntu1.name}"
                  comments = "nat rule 4"
                  enabled = true
                  method = "static"
                                  }

                              }
