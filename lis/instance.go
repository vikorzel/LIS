package lis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

type instance struct {
	endpoint  string
	userid    uint64
	username  string
	password  string
	groupname string
	groupId   uint64
	cookie    *cookiejar.Jar
	http_cli  *http.Client
	faketime  *string
}

func NewInstance(endpoint string, username string, password string, groupname string) *instance {
	inst := instance{
		endpoint:  endpoint,
		username:  username,
		password:  password,
		groupname: groupname,
	}
	cookiejar, err := cookiejar.New(
		&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	inst.cookie = cookiejar
	return &inst
}

func (inst *instance) GetEndpoint() string {
	return inst.endpoint
}

func (inst *instance) SetFaketime(new_time string) {
	faketime := string(new_time)
	inst.faketime = &faketime
}

func (inst *instance) GetFakeTime() string {
	return *(inst.faketime)
}

func (inst *instance) GetGroupId() uint64 {
	return inst.groupId
}

func (inst *instance) GetUserId() uint64 {
	return inst.userid
}

func (inst *instance) initClient() error {
	if inst.http_cli != nil {
		return nil
	}
	inst.http_cli = &http.Client{
		Jar: inst.cookie,
	}
	return nil
}

func getSessions(inst *instance) (int, error) {
	inst.initClient()
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/sessions", inst.endpoint),
		nil,
	)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := inst.http_cli.Do(req)
	if resp != nil {
		return resp.StatusCode, err
	}
	return 0, err

}

func postSessions(inst *instance) (int, *http.Response, error) {
	inst.initClient()

	credentials := SessionRequest{
		Groupname: inst.groupname,
		Username:  inst.username,
		Password:  inst.password,
	}

	bData, _ := json.Marshal(credentials)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/sessions", inst.endpoint),
		bytes.NewReader(bData),
	)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := inst.http_cli.Do(req)
	if resp != nil {
		return resp.StatusCode, resp, err
	}
	return 0, nil, err
}

func (inst *instance) Authorise() error {
	code, _ := getSessions(inst)
	if code == 403 {
		code, response, _ := postSessions(inst)
		if code != 200 {
			return fmt.Errorf(
				"wrong credentials for %s",
				inst.username,
			)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		session := Session{}
		json.Unmarshal(body, &session)
		inst.groupId = session.GroupId
		inst.userid = session.UserId

	}
	code, err := getSessions(inst)
	if code != 200 {
		return err
	}
	return nil
}
