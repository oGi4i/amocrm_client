package amocrm

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

var (
	contactArrayFields = []string{"tags", "custom_fields"}
)

func (c *ClientInfo) AddContact(contact *ContactPost) (int, error) {
	if contact.Name == "" {
		return 0, errors.New("name is empty")
	}

	url := c.Url + apiUrls["contacts"]
	resp, err := c.DoPost(url, &ContactSetRequest{Add: []*ContactPost{contact}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}

func (c *ClientInfo) GetContact(reqParams *ContactRequestParams) ([]*ContactResponse, error) {
	addValues := map[string]string{}
	contacts := new(ContactGetResponse)
	var err error

	if reqParams.ID != 0 {
		addValues["id"] = strconv.Itoa(reqParams.ID)

	} else {
		if reqParams.LimitRows != 0 {
			addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
			if reqParams.LimitOffset != 0 {
				addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
			}
		}
		if reqParams.ResponsibleUserID != 0 {
			addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
		}
		if reqParams.Query != "" {
			addValues["query"] = reqParams.Query
		}
	}

	url := c.Url + apiUrls["contacts"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, contacts)
	if err != nil {
		// fix bad json serialization, where nil array is serialized as nil object
		stringBody := string(body)
		for _, s := range contactArrayFields {
			stringBody = strings.ReplaceAll(stringBody, "\""+s+"\":{}", "\""+s+"\":[]")
		}
		err = json.Unmarshal([]byte(stringBody), contacts)
		if err != nil {
			return nil, err
		}
	}
	return contacts.Embedded.Items, err
}
