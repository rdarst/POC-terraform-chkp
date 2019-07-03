package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func dataSourceServiceUdp() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceServiceUdpRead,

                Schema: map[string]*schema.Schema{
							        	"port": {
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

func dataSourceServiceUdpRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)

	id, err := client.ReadServiceUdpData(name)
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
	readServiceUdp := chkp.ServiceUdp{}
  json.Unmarshal(id, &readServiceUdp)
	d.SetId(readServiceUdp.Uid)
  d.Set("uid", readServiceUdp.Uid)
	d.Set("name", readServiceUdp.Name)
	d.Set("port", readServiceUdp.Port)
	
	return nil
}
