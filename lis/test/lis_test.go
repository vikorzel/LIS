package lis

import (
	"LIS/lis"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSay(t *testing.T) {
	response := lis.Say()
	if response != "Hello There" {
		t.Errorf("Wrong response")
	}
}

func TestInstance(t *testing.T) {
	testsrvr := httptest.NewServer(
		http.HandlerFunc(mainHandler),
	)
	ep_test := testsrvr.URL

	defer testsrvr.Close()

	inst := lis.NewInstance(ep_test, "wrong_user", "wrong_password", "wrong_group")
	if inst.GetEndpoint() != ep_test {
		t.Errorf("Wrong Enpoint")
	}
	err := inst.Authorise()
	if err != nil && err.Error() != "wrong credentials for wrong_user" {
		t.Errorf("Authorisation is not working (1)")
	}

	err = lis.NewInstance(ep_test, "wrong_user", "TEST", "TEST").Authorise()

	if err != nil && err.Error() != "wrong credentials for wrong_user" {
		t.Errorf("Authorisation is not working (2)")
	}

	inst = lis.NewInstance(
		ep_test,
		"TEST",
		"TEST",
		"TEST",
	)
	err = inst.Authorise()
	if err != nil {
		t.Errorf("Authorisation is not working (3): %s", err.Error())
	}

	if inst.GetUserId() != 123 {
		t.Errorf("Wrong User ID: %d vs 123", inst.GetUserId())
	}

	if inst.GetGroupId() != 1234 {
		t.Errorf("Wrong Group ID: %d vs 1234", inst.GetUserId())
	}
}

func sendError(w http.ResponseWriter) {
	w.WriteHeader(403)
	resp := make(map[string]string)
	resp["error"] = "Forbidden"
	resp_text, _ := json.Marshal(resp)
	w.Write(resp_text)
}

func setSessionCookie(w *http.ResponseWriter, val string) {
	cookie := http.Cookie{}
	cookie.Name = "session"
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Value = val
	http.SetCookie(*w, &cookie)
}

func sessionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		_, err := r.Cookie("session")
		if err != nil {
			sendError(w)
			return
		}
		resp := make(map[string][]lis.Session)
		resp["sessions"] = make([]lis.Session, 1)
		resp["sessions"][0] = lis.Session{
			GroupId:   1234,
			Id:        "user-id-1",
			LastLogin: "2022-01-01 10:11:12",
			UserId:    123,
		}
		resp_text, _ := json.Marshal(resp)
		w.Write(resp_text)
	}
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			sendError(w)
			return
		}
		auth_request := lis.SessionRequest{}
		err = json.Unmarshal(body, &auth_request)
		if err != nil {
			sendError(w)
			return
		}
		if auth_request.Groupname != "TEST" || auth_request.Password != "TEST" || auth_request.Username != "TEST" {
			sendError(w)
			return
		}
		resp := lis.Session{
			GroupId:   1234,
			Id:        "user-id-1",
			LastLogin: "2022-01-01 10:11:12",
			UserId:    123,
		}
		resp_text, _ := json.Marshal(resp)
		setSessionCookie(&w, ".session1")
		w.Write(resp_text)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/sessions" {
		sessionsHandler(w, r)
	}

}
