package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)

func dataSourceSecurityZone() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceSecurityZoneRead,

                Schema: map[string]*schema.Schema{

						            "name": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "color": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Default: "black",
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                            },

                },
        }
}


func dataSourceSecurityZoneRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)
  // Call the API to get SecurityZone info
  id, err := client.ShowSecurityZone(name)
  if err != nil {
    status := err.Error()
    if (status == "404") {
          // If the object is not found remove it from state
          d.SetId("")
          return nil
    } else {
      return err
    }
  }
	readSecurityZone := chkp.SecurityZone{}
  json.Unmarshal(id, &readSecurityZone)
	d.SetId(readSecurityZone.Uid)
	d.Set("color", readSecurityZone.Color)
	d.Set("name", readSecurityZone.Name)

	return nil
}
