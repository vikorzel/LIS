package lis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type Schedule struct {
	resources []Resource
	users     []User
	timeSlots []TimeSlot
	bookings  []Boooking
	session   *instance
}

func NewSchedule(session *instance) (*Schedule, error) {
	if session == nil {
		return nil, errors.New("bad session pointer")
	}
	if session.GetUserId() == 0 {
		return nil, errors.New("not authorised session")
	}
	sch := Schedule{
		session: session,
	}
	return &sch, nil
}

func (sched *Schedule) Refresh() error {
	sched.users = sched.getUsers()
	sched.resources = sched.getResources()
	sched.timeSlots = sched.getTimeSlots()
	sched.bookings = sched.getBookings()
	return nil
}

func (sched *Schedule) getter(resname string, mapobj interface{}) error {
	_, resp, err := Get(sched.session, resname)
	log.Printf("Required %s", resname)
	if err != nil {
		log.Printf("Request %s failed: %s", resname, err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of %s request: %s", resname, err.Error())
		return err
	}
	err = json.Unmarshal(body, mapobj)
	if err != nil {
		log.Printf("Failed on %s JSON unmarshaling: %s", resname, err.Error())
		return err
	}
	return nil
}

func (sched *Schedule) getUsers() []User {
	type UserResponse struct {
		Users []User `json:"users"`
	}
	var users UserResponse
	err := sched.getter("users", &users)
	if err != nil {
		log.Fatalf("Can not get users: %s", err.Error())
	}
	return users.Users
}

func (sched *Schedule) getResources() []Resource {
	type ResourecesReponse struct {
		Resources []Resource `json:"resources"`
	}
	var resources ResourecesReponse
	err := sched.getter("resources", &resources)
	if err != nil {
		log.Fatalf("Can't get resources: %s", err.Error())
	}
	return resources.Resources
}

func (sched *Schedule) getTimeSlots() []TimeSlot {
	type TimeSlotReponse struct {
		TimeSlots []TimeSlot `json:"time_slots"`
	}
	var timeSlots TimeSlotReponse
	err := sched.getter("time_slots", &timeSlots)
	if err != nil {
		log.Fatalf("Can't get time slots: %s", err.Error())
	}
	return timeSlots.TimeSlots
}

func (sched *Schedule) getBookings() []Boooking {
	type BookingsResponse struct {
		Bookings []Boooking `json:"bookings"`
	}
	var bookings BookingsResponse
	date := sched.getDate()
	uri := fmt.Sprintf("bookings/week/%d/%02d/%02d", date.Year(), date.Month(), date.Day())
	err := sched.getter(uri, &bookings)
	if err != nil {
		log.Fatalf("Can't get bookings: %s", err.Error())
	}
	return bookings.Bookings
}

func (sched *Schedule) getDate() time.Time {
	faketime := sched.session.GetFakeTime()
	var tprocess time.Time
	if faketime == "" {
		tprocess = time.Now()
	} else {
		var err error
		tprocess, err = time.Parse("2006-01-02", faketime)
		if err != nil {
			log.Fatalf("Error with parsing fake time %s. Should have a format YYYY-MM-DD", faketime)
		}
	}
	return tprocess
}

func (sched *Schedule) GetResources() []Resource {
	return sched.resources
}
