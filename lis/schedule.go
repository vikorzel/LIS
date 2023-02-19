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

type timeTableCell struct {
	time   string
	booked bool
}

type timeTableDay struct {
	day   string
	cells []timeTableCell
}

type timeTable struct {
	name string
	id   int
	days []timeTableDay
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

func (sched *Schedule) RenderSchedule() {
	resources_cnt := 0
	for _, resource := range sched.resources {
		if resource.PrimaryFlag {
			resources_cnt++
		}
	}
	schedule := make([]timeTable, resources_cnt)
	for _, resource := range sched.resources {
		if resource.PrimaryFlag {
			schedule[resources_cnt-1].name = resource.Description
			schedule[resources_cnt-1].id = resource.ID
			resources_cnt--
		}

	}

	bookedMask := make(map[string]bool)

	for _, booking := range sched.bookings {
		mask := fmt.Sprintf("%d:%d", booking.ResourceID, booking.BookedTimeSlotID)
		bookedMask[mask] = true
	}

	for index, res := range schedule {
		schedule[index].days = make([]timeTableDay, 7)
		schedule[index].days[0].day = "Sun"
		schedule[index].days[0].cells = make([]timeTableCell, 0)

		schedule[index].days[1].day = "Mon"
		schedule[index].days[1].cells = make([]timeTableCell, 0)

		schedule[index].days[2].day = "Tue"
		schedule[index].days[2].cells = make([]timeTableCell, 0)

		schedule[index].days[3].day = "Wed"
		schedule[index].days[3].cells = make([]timeTableCell, 0)

		schedule[index].days[4].day = "Thu"
		schedule[index].days[4].cells = make([]timeTableCell, 0)

		schedule[index].days[5].day = "Fri"
		schedule[index].days[5].cells = make([]timeTableCell, 0)

		schedule[index].days[6].day = "Sat"
		schedule[index].days[6].cells = make([]timeTableCell, 0)

		for _, time_slot := range sched.timeSlots {
			mask := fmt.Sprintf("%d:%d", res.id, time_slot.ID)
			_, ok := bookedMask[mask]
			if ok {
				schedule[index].days[time_slot.DayOfWeek-1].cells = append(schedule[index].days[time_slot.DayOfWeek-1].cells, timeTableCell{
					time:   time_slot.Description,
					booked: true,
				})
			} else {
				schedule[index].days[time_slot.DayOfWeek-1].cells = append(schedule[index].days[time_slot.DayOfWeek-1].cells, timeTableCell{
					time:   time_slot.Description,
					booked: false,
				})
			}
		}
	}

	// TODO: booked_time_slot_id is not ID of time_slot, so, at first we need to request booked_time_slots and find there time_slot_id

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
