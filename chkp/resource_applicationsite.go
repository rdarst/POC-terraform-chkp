package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
        "reflect"
)


func resourceApplicationSite() *schema.Resource {
        return &schema.Resource{
                Create: resourceApplicationSiteCreate,
                Read:   resourceApplicationSiteRead,
                Update: resourceApplicationSiteUpdate,
                Delete: resourceApplicationSiteDelete,

                Schema: map[string]*schema.Schema{

						            "name": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "primarycategory": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "color": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Default: "black",
                        },
                        "urllist": &schema.Schema{
                                Type:     schema.TypeSet,
                                Required: true,
                                Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "additionalcategories": &schema.Schema{
                                Type:     schema.TypeSet,
                                Optional: true,
                                Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "tags": &schema.Schema{
                                Type:     schema.TypeSet,
                                Optional: true,
                                Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "description": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "urlsdefinedasregularexpression": {
                                Type:     schema.TypeBool,
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

func resourceApplicationSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var applicationsite = chkp.ApplicationSite{}
	applicationsite.Name = d.Get("name").(string)
  applicationsite.Color = d.Get("color").(string)
  applicationsite.PrimaryCategory = d.Get("primarycategory").(string)
  applicationsite.Description = d.Get("description").(string)
  applicationsite.Urlsdefinedasregularexpression = d.Get("urlsdefinedasregularexpression").(bool)

  // Pull in the URL list
  urllist := d.Get("urllist").(*schema.Set).List()
  urllist_result := make([]string, 0, len(urllist))
	for _, v := range urllist {
		val, ok := v.(string)
		if ok && val != "" {
			urllist_result = append(urllist_result, v.(string))
		}
  }
  applicationsite.URLList = urllist_result

  // Pull in the additional categories
  additionalcategories := d.Get("additionalcategories").(*schema.Set).List()
  additionalcategories_result := make([]string, 0, len(additionalcategories))
  for _, v := range additionalcategories {
    val, ok := v.(string)
    if ok && val != "" {
      additionalcategories_result = append(additionalcategories_result, v.(string))
    }
  }
  applicationsite.AdditionalCategories = additionalcategories_result

  // Pull in the tags
  taglist := d.Get("tags").(*schema.Set).List()
  taglist_result := make([]string, 0, len(taglist))
	for _, v := range taglist {
		val, ok := v.(string)
		if ok && val != "" {
			taglist_result = append(taglist_result, v.(string))
		}
  }
  applicationsite.Tags = taglist_result

  id, err := client.CreateApplicationSite(applicationsite)
  readApplicationSite := chkp.ApplicationSite{}
  json.Unmarshal(id, &readApplicationSite)
	d.SetId(readApplicationSite.Uid)
  d.Set("uid", readApplicationSite.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceApplicationSiteRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
  // Call the API to get ApplicationSite info
  id, err := client.ShowApplicationSite(d.Id())
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
	readApplicationSite := chkp.ApplicationSiteResult{}
  json.Unmarshal(id, &readApplicationSite)
	d.SetId(readApplicationSite.Uid)
	d.Set("color", readApplicationSite.Color)
	d.Set("name", readApplicationSite.Name)
  d.Set("primarycategory", readApplicationSite.PrimaryCategory)
  d.Set("description", readApplicationSite.Description)
  d.Set("urlsdefinedasregularexpression", readApplicationSite.Urlsdefinedasregularexpression)

  // Pull in the list of  URLs

  urllist := make([]string, 0, len(readApplicationSite.URLList))
  for _, url := range readApplicationSite.URLList {
    urllist = append(urllist, url)
  }
  d.Set("urllist", client.ConvertListtoSet(urllist))

  // Pull in the list of additonal categories
  additionalcategorieslist := make([]string, 0, len(readApplicationSite.AdditionalCategories))
  for _, additionalcategories := range readApplicationSite.AdditionalCategories {
    additionalcategorieslist = append(additionalcategorieslist, additionalcategories)
  }
  d.Set("additionalcategories", client.ConvertListtoSet(additionalcategorieslist))


  // Pull in the list of Tags

  i := 0
  tag := make([]string, 0, len(readApplicationSite.Tags))
  for _, taglist := range readApplicationSite.Tags {
        tag = append(tag, taglist.Name)
        i += 1
     }

  d.Set("tags", tag)

	return nil
}

func resourceApplicationSiteUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var applicationsite = chkp.ApplicationSite{}
  applicationsite.Color = d.Get("color").(string)
  applicationsite.PrimaryCategory = d.Get("primarycategory").(string)
  applicationsite.Description = d.Get("description").(string)
  applicationsite.Urlsdefinedasregularexpression = d.Get("urlsdefinedasregularexpression").(bool)

  if d.HasChange("name") {applicationsite.Newname = d.Get("name").(string)}
	applicationsite.Uid = d.Get("uid").(string)

  // Pull in the URL list
  urllist := d.Get("urllist").(*schema.Set).List()
  urllist_result := make([]string, 0, len(urllist))
  for _, v := range urllist {
    val, ok := v.(string)
    if ok && val != "" {
      urllist_result = append(urllist_result, v.(string))
    }
  }
  applicationsite.URLList = urllist_result

  // Pull in the additional categories
  additionalcategories := d.Get("additionalcategories").(*schema.Set).List()
  additionalcategories_result := make([]string, 0, len(additionalcategories))
  for _, v := range additionalcategories {
    val, ok := v.(string)
    if ok && val != "" {
      additionalcategories_result = append(additionalcategories_result, v.(string))
    }
  }
  applicationsite.AdditionalCategories = additionalcategories_result

  o, n := d.GetChange("tags")
  	if o == nil {
  			o = new(schema.Set)
  	}
  		if n == nil {
  			n = new(schema.Set)
  		}

  otaglist := o.(*schema.Set).List()
  ntaglist := n.(*schema.Set).List()
  if !(reflect.DeepEqual(otaglist, ntaglist)) {
    // Pull in the tags if there are changes
    // If the tags list is empty remove the old tags using the remove api call
    // If not empty fill the array with tags to be set
    if len(ntaglist) > 0 {
    taglist := d.Get("tags").(*schema.Set).List()
    taglist_result := make([]string, 0, len(taglist))
    for _, v := range taglist {
      val, ok := v.(string)
      if ok && val != "" {
        taglist_result = append(taglist_result, v.(string))
      }
    }
    applicationsite.Tags = taglist_result
  }else {

    taglist_result := make([]string, 0, len(otaglist))
    for _, v := range otaglist {
      val, ok := v.(string)
      if ok && val != "" {
        taglist_result = append(taglist_result, v.(string))
      }
    }
    // Remove the old Tags
    var applicationsitetagupdate = chkp.ApplicationSiteTagAddRemove{}
    applicationsitetagupdatetags := chkp.TagsAddRemove{
      Remove:       taglist_result,
    }
    applicationsitetagupdate.Uid = d.Get("uid").(string)
    applicationsitetagupdate.TagsAddRemove = applicationsitetagupdatetags
    id, err := client.SetApplicationSiteUpdateTag(applicationsitetagupdate)
    _ = id
    if err != nil {
  		return err
  	}
  }
}

	id, err := client.SetApplicationSite(applicationsite)
  readApplicationSite := chkp.ApplicationSiteResult{}
  if err := json.Unmarshal(id, &readApplicationSite); err != nil {return err}
  //json.Unmarshal(id, &readApplicationSite)
	d.SetId(readApplicationSite.Uid)
  d.Set("uid", readApplicationSite.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceApplicationSiteDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  used, err := client.CheckWhereUsed(d.Id())
  if used > 0 {
    client.WaitUntilNotUsed(d.Id())
  }
	result, err := client.DeleteApplicationSite(d.Id())
  _ = result
  if err != nil {
		return err
	}
  return nil
}
