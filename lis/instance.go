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
	userID    uint64
	username  string
	password  string
	groupname string
	groupID   uint64
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
		userID:    0,
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
	if inst.faketime != nil {
		return *(inst.faketime)
	}
	return ""
}

func (inst *instance) GetGroupId() uint64 {
	return inst.groupID
}

func (inst *instance) GetUserId() uint64 {
	return inst.userID
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

func Get(inst *instance, handler string) (int, *http.Response, error) {
	log.Printf("try to get %s", handler)
	inst.initClient()
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", inst.endpoint, handler),
		nil,
	)
	if err != nil {
		log.Printf("error with request: %s", err.Error())
		return 0, nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := inst.http_cli.Do(req)
	if resp != nil {
		log.Printf("Received response (%d)", resp.StatusCode)
		return resp.StatusCode, resp, err
	}
	log.Printf("Data sending failed: %s", err.Error())
	return 0, nil, err

}

func Post(inst *instance, handler string, payload *[]byte) (int, *http.Response, error) {
	log.Printf("Try to post %s handler", handler)
	inst.initClient()
	var req *http.Request
	var err error
	var reader *bytes.Reader = nil
	if payload != nil {
		reader = bytes.NewReader(*payload)
	}
	req, err = http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", inst.endpoint, handler),
		reader,
	)
	if err != nil {
		log.Printf("error with request building: %s", err.Error())
		return 0, nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := inst.http_cli.Do(req)
	if err != nil {
		log.Printf("error with request sending: %s", err.Error())
		return 0, nil, err
	}
	return resp.StatusCode, resp, err

}

func getSessions(inst *instance) (int, error) {
	code, _, err := Get(inst, "sessions")
	return code, err
}

func postSessions(inst *instance) (int, *http.Response, error) {
	inst.initClient()
	log.Println("Post session try")
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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
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
			buf := make([]byte, response.ContentLength+1)
			response.Body.Read(buf)
			return fmt.Errorf(
				"wrong credentials for %s: %s",
				inst.username,
				string(buf[:]),
			)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		session := SessionResonse{}
		json.Unmarshal(body, &session)
		inst.groupID = session.GroupID
		inst.userID = session.UserID

	}
	code, err := getSessions(inst)
	if code != 200 {
		return err
	}
	return nil
}
