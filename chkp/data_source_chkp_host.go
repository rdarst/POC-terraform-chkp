package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func dataSourceHost() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceHostRead,

                Schema: map[string]*schema.Schema{
							        	"ipv4address": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
						            "name": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
						            "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                },
        }
}


func dataSourceHostRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)

	id, err := client.ReadHostData(name)
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
	readHost := chkp.Host{}
  json.Unmarshal(id, &readHost)
	d.SetId(readHost.Uid)
  d.Set("uid", readHost.Uid)
	d.Set("name", readHost.Name)
	d.Set("ipv4address", readHost.Ipv4address)
	
	return nil
}
