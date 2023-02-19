package lis

import (
	"LIS/lis"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestSchedule(t *testing.T) {
	testsrvr := httptest.NewServer(
		http.HandlerFunc(mainHandler),
	)
	instance := lis.NewInstance(
		testsrvr.URL,
		"TEST",
		"TEST",
		"TEST",
	)

	_, err := lis.NewSchedule(nil)
	if err == nil || err.Error() != "bad session pointer" {
		t.Errorf("Session pointer for schedule is not checked")
	}

	_, err = lis.NewSchedule(instance)
	if err == nil || err.Error() != "not authorised session" {
		t.Errorf("Session auth is not checked for the schedule credentials")
	}

	err = instance.Authorise()
	if err != nil {
		t.Error("Auth credentials is not valid for the end user")
	}

	sched, err := lis.NewSchedule(instance)
	if err != nil {
		t.Errorf("Can not create schedule obj with err: %s", err.Error())
	}

	instance.SetFaketime("2022-11-29")
	sched.Refresh()
	resources := sched.GetResources()
	if len(resources) < 1 {
		t.Errorf("Didn't receieved resources information")
	}
	sched.RenderSchedule()
	t.Errorf("The end")

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

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(
		[]byte("{\n  \"resources\": [\n    {\n      \"description\": \"Instructor John Doe\", \n      \"group_id\": 19618, \n      \"id\": 77790, \n      \"primary_flag\": false, \n      \"redact_booking_text\": false, \n      \"sequence_num\": 1, \n      \"symbol\": \"diapur\"\n    }, \n    {\n      \"description\": \"Cessna 172\", \n      \"group_id\": 19618, \n      \"id\": 77787, \n      \"primary_flag\": true, \n      \"redact_booking_text\": false, \n      \"sequence_num\": 1, \n      \"symbol\": \"none\"\n    }, \n    {\n      \"description\": \"Instructor Fred Bloggs\", \n      \"group_id\": 19618, \n      \"id\": 77789, \n      \"primary_flag\": false, \n      \"redact_booking_text\": false, \n      \"sequence_num\": 2, \n      \"symbol\": \"diayel\"\n    }, \n    {\n      \"description\": \"Piper Archer\", \n      \"group_id\": 19618, \n      \"id\": 77791, \n      \"primary_flag\": true, \n      \"redact_booking_text\": false, \n      \"sequence_num\": 2, \n      \"symbol\": \"none\"\n    }, \n    {\n      \"description\": \"Club Life Raft\", \n      \"group_id\": 19618, \n      \"id\": 77788, \n      \"primary_flag\": false, \n      \"redact_booking_text\": false, \n      \"sequence_num\": 3, \n      \"symbol\": \"square_light_green\"\n    }\n  ]\n}\n"),
	)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(
		[]byte("{\n  \"users\": [\n    {\n      \"administrator\": true, \n      \"allow_alter_others\": true, \n      \"allow_edits\": true, \n      \"booking_change_emails\": false, \n      \"confirmation_email_sent\": false, \n      \"email\": \"admin@e-allocator.com\", \n      \"email_confirmed\": false, \n      \"email_messages\": false, \n      \"email_preferences\": true, \n      \"group_id\": 19618, \n      \"id\": 359235, \n      \"last_login\": \"2007-12-20 03:49:00\", \n      \"member_details_private\": false, \n      \"name\": \"Demo Administrator\", \n      \"username\": \"ADMIN\"\n    }, \n    {\n      \"administrator\": false, \n      \"allow_alter_others\": true, \n      \"allow_edits\": true, \n      \"booking_change_emails\": true, \n      \"confirmation_email_sent\": false, \n      \"email\": \"demo@e-allocator.com\", \n      \"email_confirmed\": true, \n      \"email_messages\": true, \n      \"email_preferences\": true, \n      \"group_id\": 19618, \n      \"id\": 360847, \n      \"last_login\": \"2022-11-29 19:37:32\", \n      \"member_details_private\": false, \n      \"name\": \"Demo User\", \n      \"username\": \"DEMO\"\n    }\n  ]\n}\n"),
	)
}

func timeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(
		[]byte("{\n  \"time_slots\": [\n    {\n      \"day_of_week\": 2, \n      \"description\": \"9am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759160, \n      \"prime\": false, \n      \"sequence_num\": 1\n    }, \n    {\n      \"day_of_week\": 2, \n      \"description\": \"2pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759161, \n      \"prime\": false, \n      \"sequence_num\": 2\n    }, \n    {\n      \"day_of_week\": 3, \n      \"description\": \"9am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759162, \n      \"prime\": false, \n      \"sequence_num\": 3\n    }, \n    {\n      \"day_of_week\": 3, \n      \"description\": \"2pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759163, \n      \"prime\": false, \n      \"sequence_num\": 4\n    }, \n    {\n      \"day_of_week\": 4, \n      \"description\": \"9am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759164, \n      \"prime\": false, \n      \"sequence_num\": 5\n    }, \n    {\n      \"day_of_week\": 4, \n      \"description\": \"2pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759165, \n      \"prime\": false, \n      \"sequence_num\": 6\n    }, \n    {\n      \"day_of_week\": 5, \n      \"description\": \"9am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759166, \n      \"prime\": false, \n      \"sequence_num\": 7\n    }, \n    {\n      \"day_of_week\": 5, \n      \"description\": \"2pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759167, \n      \"prime\": false, \n      \"sequence_num\": 8\n    }, \n    {\n      \"day_of_week\": 6, \n      \"description\": \"9am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759168, \n      \"prime\": false, \n      \"sequence_num\": 9\n    }, \n    {\n      \"day_of_week\": 6, \n      \"description\": \"2pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759169, \n      \"prime\": false, \n      \"sequence_num\": 10\n    }, \n    {\n      \"day_of_week\": 7, \n      \"description\": \"9am - 11:30pm\", \n      \"group_id\": 19618, \n      \"id\": 759170, \n      \"prime\": false, \n      \"sequence_num\": 11\n    }, \n    {\n      \"day_of_week\": 7, \n      \"description\": \"11:30am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759171, \n      \"prime\": false, \n      \"sequence_num\": 12\n    }, \n    {\n      \"day_of_week\": 7, \n      \"description\": \"2pm - 4:30pm\", \n      \"group_id\": 19618, \n      \"id\": 759172, \n      \"prime\": false, \n      \"sequence_num\": 13\n    }, \n    {\n      \"day_of_week\": 7, \n      \"description\": \"4:30pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759173, \n      \"prime\": false, \n      \"sequence_num\": 14\n    }, \n    {\n      \"day_of_week\": 1, \n      \"description\": \"9am - 11:30pm\", \n      \"group_id\": 19618, \n      \"id\": 759174, \n      \"prime\": false, \n      \"sequence_num\": 15\n    }, \n    {\n      \"day_of_week\": 1, \n      \"description\": \"11:30am - 2pm\", \n      \"group_id\": 19618, \n      \"id\": 759175, \n      \"prime\": false, \n      \"sequence_num\": 16\n    }, \n    {\n      \"day_of_week\": 1, \n      \"description\": \"2pm - 4:30pm\", \n      \"group_id\": 19618, \n      \"id\": 759176, \n      \"prime\": false, \n      \"sequence_num\": 17\n    }, \n    {\n      \"day_of_week\": 1, \n      \"description\": \"4:30pm - 7pm\", \n      \"group_id\": 19618, \n      \"id\": 759177, \n      \"prime\": false, \n      \"sequence_num\": 18\n    }\n  ]\n}\n"),
	)
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(
		[]byte("{\n  \"bookings\": [\n    {\n      \"block_uuid\": null, \n      \"booked_by_user_id\": 360847, \n      \"booked_time_slot_id\": 7805732, \n      \"booked_when\": \"2022-11-29 00:00:00\", \n      \"description\": \"1234\", \n      \"id\": 11764275, \n      \"primary_booking_id\": null, \n      \"resource_id\": 77787, \n      \"secondary_bookings\": []\n    }\n  ]\n}\n"),
	)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/sessions" {
		sessionsHandler(w, r)
	} else if r.RequestURI == "/resources" {
		resourcesHandler(w, r)
	} else if r.RequestURI == "/users" {
		usersHandler(w, r)
	} else if r.RequestURI == "/time_slots" {
		timeSlotsHandler(w, r)
	} else if r.RequestURI == "/bookings/week/2022/11/29" {
		bookingHandler(w, r)
	} else {
		w.WriteHeader(404)
	}

}
