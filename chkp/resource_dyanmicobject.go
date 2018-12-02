package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceDynamicObject() *schema.Resource {
        return &schema.Resource{
                Create: resourceDynamicObjectCreate,
                Read:   resourceDynamicObjectRead,
                Update: resourceDynamicObjectUpdate,
                Delete: resourceDynamicObjectDelete,

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

func resourceDynamicObjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var dynamicobject = chkp.DynamicObject{}
	dynamicobject.Name = d.Get("name").(string)
  dynamicobject.Color = d.Get("color").(string)

  id, err := client.CreateDynamicObject(dynamicobject)
  readDynamicObject := chkp.DynamicObject{}
  json.Unmarshal(id, &readDynamicObject)
	d.SetId(readDynamicObject.Uid)
  d.Set("uid", readDynamicObject.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceDynamicObjectRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
  // Call the API to get DynamicObject info
  id, err := client.ShowDynamicObject(d.Id())

	readDynamicObject := chkp.DynamicObject{}
  json.Unmarshal(id, &readDynamicObject)
	d.SetId(readDynamicObject.Uid)
	d.Set("color", readDynamicObject.Color)
	d.Set("name", readDynamicObject.Name)

  if err != nil {
		return err
	}
	return nil
}

func resourceDynamicObjectUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var dynamicobject = chkp.DynamicObject{}
	dynamicobject.Color = d.Get("color").(string)
  if d.HasChange("name") {dynamicobject.Newname = d.Get("name").(string)}
	dynamicobject.Uid = d.Get("uid").(string)

	id, err := client.SetDynamicObject(dynamicobject)

  readDynamicObject := chkp.DynamicObject{}
  if err := json.Unmarshal(id, &readDynamicObject); err != nil {return err}
  //json.Unmarshal(id, &readDynamicObject)
	d.SetId(readDynamicObject.Uid)
  d.Set("uid", readDynamicObject.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceDynamicObjectDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteDynamicObject(d.Id())
	return nil

}
