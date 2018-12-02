package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceAccessLayer() *schema.Resource {
        return &schema.Resource{
                Create: resourceAccessLayerCreate,
                Read:   resourceAccessLayerRead,
                Update: resourceAccessLayerUpdate,
                Delete: resourceAccessLayerDelete,

                Schema: map[string]*schema.Schema{
							        	"appandurl": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                  },
                        "contentawareness": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: false,
                        },
                        "adddefaultrule": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: true,
                        },
                        "firewall": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: true,
                        },
                        "mobileaccess": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: false,
                        },
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
                        "shared": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: false,
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                },
        }
}

func resourceAccessLayerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var accesslayer = chkp.AccessLayer{}
	accesslayer.Name = d.Get("name").(string)
	accesslayer.AppAndUrl = d.Get("appandurl").(bool)
  accesslayer.ContentAwareness = d.Get("contentawareness").(bool)
  accesslayer.Firewall = d.Get("firewall").(bool)
  accesslayer.MobileAccess = d.Get("mobileaccess").(bool)
  accesslayer.Shared = d.Get("shared").(bool)
  accesslayer.AddDefaultRule = d.Get("adddefaultrule").(bool)
  accesslayer.Color = d.Get("color").(string)
  accesslayer.Comments = d.Get("comments").(string)
	id, err := client.CreateAccessLayer(accesslayer)
  //Read in from the API Output
  readAccessLayer := chkp.AccessLayer{}
  json.Unmarshal(id, &readAccessLayer)
	d.SetId(readAccessLayer.Uid)
  d.Set("uid", readAccessLayer.Uid)
  d.Set("name", readAccessLayer.Name)
	if err != nil {
		return err
	}
	return nil
}

func resourceAccessLayerRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
	id, err := client.ShowAccessLayer(d.Id())

	readAccessLayer := chkp.AccessLayer{}
  json.Unmarshal(id, &readAccessLayer)
	d.SetId(readAccessLayer.Uid)
	d.Set("color", readAccessLayer.Color)
	d.Set("name", readAccessLayer.Name)
	d.Set("appandurl", readAccessLayer.AppAndUrl)
  d.Set("contentawareness", readAccessLayer.ContentAwareness)
  d.Set("firewall", readAccessLayer.Firewall)
  d.Set("mobileaccess", readAccessLayer.MobileAccess)
  d.Set("shared", readAccessLayer.Shared)
  d.Set("comments", readAccessLayer.Comments)
  if err != nil {
		return err
	}
	return nil
}

func resourceAccessLayerUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var accesslayer = chkp.AccessLayer{}
  accesslayer.Name = d.Get("name").(string)
	accesslayer.AppAndUrl = d.Get("appandurl").(bool)
  accesslayer.ContentAwareness = d.Get("contentawareness").(bool)
  accesslayer.Firewall = d.Get("firewall").(bool)
  accesslayer.MobileAccess = d.Get("mobileaccess").(bool)
  accesslayer.Shared = d.Get("shared").(bool)
  accesslayer.Color = d.Get("color").(string)
  accesslayer.Comments = d.Get("comments").(string)
	id, err := client.SetAccessLayer(accesslayer)
  if err != nil {return err}
  //Read in from the API Output
  readAccessLayer := chkp.AccessLayer{}
  if err := json.Unmarshal(id, &readAccessLayer); err != nil {return err}
  d.SetId(readAccessLayer.Uid)
  d.Set("uid", readAccessLayer.Uid)
  d.Set("name", readAccessLayer.Name)
	if err != nil {
		return err
	}
	return nil
}

func resourceAccessLayerDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteAccessLayer(d.Id())
	return nil

}
