package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceTag() *schema.Resource {
        return &schema.Resource{
                Create: resourceTagCreate,
                Read:   resourceTagRead,
                Update: resourceTagUpdate,
                Delete: resourceTagDelete,

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
                        "comments": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                            },

                },
        }
}

func resourceTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var tag = chkp.Tag{}
	tag.Name = d.Get("name").(string)
  tag.Color = d.Get("color").(string)

  id, err := client.CreateTag(tag)
  readTag := chkp.Tag{}
  json.Unmarshal(id, &readTag)
	d.SetId(readTag.Uid)
  d.Set("uid", readTag.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceTagRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  // Call the API to get Tag info
  id, err := client.ShowTag(d.Id())

	readTag := chkp.Tag{}
  json.Unmarshal(id, &readTag)
	d.SetId(readTag.Uid)
	d.Set("color", readTag.Color)
	d.Set("name", readTag.Name)

  if err != nil {
		return err
	}
	return nil
}

func resourceTagUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var tag = chkp.Tag{}
	tag.Color = d.Get("color").(string)
  if d.HasChange("name") {tag.Newname = d.Get("name").(string)}
	tag.Uid = d.Get("uid").(string)

	id, err := client.SetTag(tag)

  readTag := chkp.Tag{}
  if err := json.Unmarshal(id, &readTag); err != nil {return err}
  //json.Unmarshal(id, &readTag)
	d.SetId(readTag.Uid)
  d.Set("uid", readTag.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceTagDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteTag(d.Id())
	return nil

}
