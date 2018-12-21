resource "chkp_servicetcp" "healthcheck" {
      port = "8117"
      name = "AzureHealthCheck"
      color = "blue"
      matchbysig = false
              }
