package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func dataSourceDynamicObject() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceDynamicObjectRead,
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

func dataSourceDynamicObjectRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)
  // Call the API to get DynamicObject info
  id, err := client.ShowDynamicObject(name)
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
	readDynamicObject := chkp.DynamicObject{}
  json.Unmarshal(id, &readDynamicObject)
	d.SetId(readDynamicObject.Uid)
	d.Set("color", readDynamicObject.Color)
	d.Set("name", readDynamicObject.Name)

	return nil
}
