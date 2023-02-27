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
	resources         []Resource
	users             []User
	timeSlots         []TimeSlot
	bookings          []Boooking
	booked_time_slots []BookedTimeSlot
	session           *instance
	bts2ts            map[int]int
	renderedData      []TimeTable
}

type TimeTableCell struct {
	Time   string
	Booked bool
	ID     int
}

type TimeTableDay struct {
	Day   string
	Cells []TimeTableCell
}

type TimeTable struct {
	Name string
	ID   int
	Days []TimeTableDay
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

func (sched *Schedule) makeBTS2TSMap() {
	sched.bts2ts = make(map[int]int)
	for _, booked_time_slot := range sched.booked_time_slots {
		sched.bts2ts[booked_time_slot.ID] = booked_time_slot.TimeSlotID
	}
}

func (sched *Schedule) Refresh() error {
	sched.users = sched.getUsers()
	sched.resources = sched.getResources()
	sched.timeSlots = sched.getTimeSlots()
	sched.bookings = sched.getBookings()
	sched.booked_time_slots = sched.getBookedTimeSlots()

	sched.makeBTS2TSMap()
	return nil
}

func (sched *Schedule) RenderSchedule() []TimeTable {
	resources_cnt := 0
	for _, resource := range sched.resources {
		if resource.PrimaryFlag {
			resources_cnt++
		}
	}
	schedule := make([]TimeTable, resources_cnt)
	for _, resource := range sched.resources {
		if resource.PrimaryFlag {
			schedule[resources_cnt-1].Name = resource.Description
			schedule[resources_cnt-1].ID = resource.ID
			resources_cnt--
		}

	}

	bookedMask := make(map[string]bool)

	for _, booking := range sched.bookings {
		timeSlotID := sched.bts2ts[booking.BookedTimeSlotID]
		mask := fmt.Sprintf("%d:%d", booking.ResourceID, timeSlotID)
		bookedMask[mask] = true
	}

	for index, res := range schedule {
		schedule[index].Days = make([]TimeTableDay, 7)
		schedule[index].Days[0].Day = "Sun"
		schedule[index].Days[0].Cells = make([]TimeTableCell, 0)

		schedule[index].Days[1].Day = "Mon"
		schedule[index].Days[1].Cells = make([]TimeTableCell, 0)

		schedule[index].Days[2].Day = "Tue"
		schedule[index].Days[2].Cells = make([]TimeTableCell, 0)

		schedule[index].Days[3].Day = "Wed"
		schedule[index].Days[3].Cells = make([]TimeTableCell, 0)

		schedule[index].Days[4].Day = "Thu"
		schedule[index].Days[4].Cells = make([]TimeTableCell, 0)

		schedule[index].Days[5].Day = "Fri"
		schedule[index].Days[5].Cells = make([]TimeTableCell, 0)

		schedule[index].Days[6].Day = "Sat"
		schedule[index].Days[6].Cells = make([]TimeTableCell, 0)

		for _, time_slot := range sched.timeSlots {
			mask := fmt.Sprintf("%d:%d", res.ID, time_slot.ID)
			_, ok := bookedMask[mask]
			if ok {
				schedule[index].Days[time_slot.DayOfWeek-1].Cells = append(schedule[index].Days[time_slot.DayOfWeek-1].Cells, TimeTableCell{
					Time:   time_slot.Description,
					Booked: true,
					ID:     time_slot.ID,
				})
			} else {
				schedule[index].Days[time_slot.DayOfWeek-1].Cells = append(schedule[index].Days[time_slot.DayOfWeek-1].Cells, TimeTableCell{
					Time:   time_slot.Description,
					Booked: false,
					ID:     time_slot.ID,
				})
			}
		}
	}

	sched.renderedData = schedule
	return schedule
	// TODO: booked_time_slot_id is not ID of time_slot, so, at first we need to request booked_time_slots and find there time_slot_id
}

func (sched *Schedule) bookTimeSlot(timeSlotID int, resourceID int) bool {
	// POST https://www.e-allocator.com/api/v1/booked_time_slots {time_slot_id: 759170, booking_date: "2022-12-03"} <- {"booking_date": "2022-12-03", "group_id": 19618, "id": 7805732, "time_slot_id": 759170}
	date := sched.session.getDate()
	dateFormatedString := fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())

	request := BookingTimeSlotRequest{
		TimeSlotID:  timeSlotID,
		BookingDate: dateFormatedString,
	}

	// POST https://www.e-allocator.com/api/v1/bookings {"resource_id":77787,"description":"1234","booked_time_slot_id":7805732,"booked_by_user_id":360847,"booked_when":"2022-11-29","secondary_resource_ids":[],"ical":false}
}

func (sched *Schedule) BookIfPossible(day string, time string) *string {
	if sched.renderedData == nil {
		sched.RenderSchedule()
	}
	for _, resource := range sched.renderedData {
		for _, dayCell := range resource.Days {
			if dayCell.Day == day {
				for _, timeCell := range dayCell.Cells {
					if timeCell.Time == time && timeCell.Booked == false {
						sched.bookTimeSlot(timeCell.ID, resource.ID)
					}
				}
			}
		}
	}
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

func (sched *Schedule) getBookedTimeSlots() []BookedTimeSlot {
	type BookedTimeSlotResponse struct {
		BookedTimeSlots []BookedTimeSlot `json:"booked_time_slots"`
	}
	var timeSlots BookedTimeSlotResponse
	date := sched.getDate()
	uri := fmt.Sprintf("booked_time_slots/week/%d/%02d/%02d", date.Year(), date.Month(), date.Day())
	err := sched.getter(uri, &timeSlots)
	if err != nil {
		log.Fatalf("Can't get booked time slots: %s", err.Error())
	}
	return timeSlots.BookedTimeSlots
}

func (sched *Schedule) GetResources() []Resource {
	return sched.resources
}
