package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
        "strings"
)


func resourceServiceTcp() *schema.Resource {
        return &schema.Resource{
                Create: resourceServiceTcpCreate,
                Read:   resourceServiceTcpRead,
                Update: resourceServiceTcpUpdate,
                Delete: resourceServiceTcpDelete,

                Schema: map[string]*schema.Schema{
							        	"port": {
                                Type:     schema.TypeString,
                                Required: true,
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
                        "protocol": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "sessiontimeout": {
                                Type:     schema.TypeInt,
                                Optional: true,
                                Default: 3600,
                        },
                        "matchbysig": {
                                Type:     schema.TypeBool,
                                Optional: true,
                        },
                        "matchforany": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: true,
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                },
        }
}

func resourceServiceTcpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var servicetcp = chkp.ServiceTcp{}
	servicetcp.Name = d.Get("name").(string)
	servicetcp.Port = d.Get("port").(string)
  servicetcp.Color = d.Get("color").(string)
  servicetcp.Protocol = strings.ToUpper(d.Get("protocol").(string))
  servicetcp.SessionTimeout = d.Get("sessiontimeout").(int)
  servicetcp.MatchBySig = d.Get("matchbysig").(bool)
  servicetcp.MatchForAny = d.Get("matchforany").(bool)
	id, err := client.CreateServiceTcp(servicetcp)
  //Read in from the API Output
  readServiceTcp := chkp.ServiceTcp{}
  json.Unmarshal(id, &readServiceTcp)
	d.SetId(readServiceTcp.Uid)
  d.Set("uid", readServiceTcp.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceServiceTcpRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
	id, err := client.ShowServiceTcp(d.Id())

	readServiceTcp := chkp.ServiceTcp{}
  json.Unmarshal(id, &readServiceTcp)
	d.SetId(readServiceTcp.Uid)
	d.Set("color", readServiceTcp.Color)
	d.Set("name", readServiceTcp.Name)
	d.Set("port", readServiceTcp.Port)
  d.Set("protocol", strings.ToUpper(readServiceTcp.Protocol))
  d.Set("sessiontimeout", readServiceTcp.SessionTimeout)
  d.Set("matchbysig", readServiceTcp.MatchBySig)
  d.Set("matchforany", readServiceTcp.MatchForAny)
	if err != nil {
		return err
	}
	return nil
}

func resourceServiceTcpUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var servicetcp = chkp.ServiceTcp{}
	servicetcp.Port = d.Get("port").(string)
	servicetcp.Color = d.Get("color").(string)
  servicetcp.Protocol = d.Get("protocol").(string)
  servicetcp.SessionTimeout = d.Get("sessiontimeout").(int)
  servicetcp.MatchBySig = d.Get("matchbysig").(bool)
  servicetcp.MatchForAny = d.Get("matchforany").(bool)
	if d.HasChange("name") {servicetcp.Newname = d.Get("name").(string)}
	servicetcp.Uid = d.Get("uid").(string)
	id, err := client.SetServiceTcp(servicetcp)
  if err != nil {return err}
  //Read in from the API Output
  readServiceTcp := chkp.ServiceTcp{}
  if err := json.Unmarshal(id, &readServiceTcp); err != nil {return err}
  d.SetId(readServiceTcp.Uid)
  d.Set("uid", readServiceTcp.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourceServiceTcpDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteServiceTcp(d.Id())
	return nil

}
