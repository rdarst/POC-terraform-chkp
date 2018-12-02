package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)

func resourceAccessSection() *schema.Resource {
        return &schema.Resource{
                Create: resourceAccessSectionCreate,
                Read:   resourceAccessSectionRead,
                Update: resourceAccessSectionUpdate,
                Delete: resourceAccessSectionDelete,

                Schema: map[string]*schema.Schema{

						            "name": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "position": {
                                Type:     schema.TypeList,
                                Optional: true,
                                Elem: &schema.Resource{
                                  Schema: map[string]*schema.Schema{
                                      "top": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "bottom": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "above": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "below": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                    },
                                  },
                                },
						            "layer": {
                                Type:     schema.TypeString,
                                Required: true,
                                  },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                                  },
                },
        }
}

func resourceAccessSectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var accesssection = chkp.AccessSection{}
	accesssection.Name = d.Get("name").(string)
  accesssection.Layer = d.Get("layer").(string)

  positionlist := d.Get("position").([]interface{})
  for i := range positionlist {
        positionelements := positionlist[i].(map[string]interface{})
        positionreturn := chkp.Position{
          Above:       positionelements["above"].(string),
          Below:       positionelements["below"].(string),
          Top:         positionelements["top"].(string),
          Bottom:      positionelements["bottom"].(string),
        }
        accesssection.Position = positionreturn
    }

  id, err := client.CreateAccessSection(accesssection)
  //Read in from the API Output
  readAccessSection := chkp.AccessSectionResult{}
  json.Unmarshal(id, &readAccessSection)
	d.SetId(readAccessSection.Uid)
  d.Set("uid", readAccessSection.Uid)
  if err != nil {
    return err
  }
	return nil
}

func resourceAccessSectionRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
  layer := d.Get("layer").(string)
  id, err := client.ShowAccessSection(d.Id(),layer)

	readAccessSection := chkp.AccessSectionResult{}
  json.Unmarshal(id, &readAccessSection)
	d.SetId(readAccessSection.Uid)
	d.Set("uid", readAccessSection.Uid)
  d.Set("name", readAccessSection.Name)
  if len(readAccessSection.Layer) != 0 {
    layername, err := client.ReadLayerUIDtoName(readAccessSection.Layer)
    if err == nil {
      d.Set("layer", layername)
    }
  }
  if err != nil {
		return err
	}
	return nil
}

func resourceAccessSectionUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var accesssection = chkp.AccessSectionUpdate{}
  accesssection.Layer = d.Get("layer").(string)

  //Update name if it has changed
	if d.HasChange("name") {accesssection.Newname = d.Get("name").(string)}

	accesssection.Uid = d.Get("uid").(string)

  id,  err := client.SetAccessSection(accesssection)
  if err != nil {return err}
  _ = id
  //Read in from the API Output
  //readAccessSection := chkp.AccessSection{}
  //if err := json.Unmarshal(id, &readAccessSection); err != nil {return err}
  //d.SetId(readAccessSection.Uid)
  //d.Set("uid", readAccessSection.Uid)

	//if err != nil {
	//	return err
	//}
	return nil
}

func resourceAccessSectionDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
    layer := d.Get("layer").(string)
    uid := d.Get("uid").(string)
	client.DeleteAccessSection(uid, layer)
	return nil

}
