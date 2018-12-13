resource "chkp_servicegroup" "myservicegroup" {
      name = "myservicegroup"
      color = "blue"
      members = ["dns","http","https"]
      }
