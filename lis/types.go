package lis

type Session struct {
	GroupId   uint64  `json:"group_id"`
	Id        string  `json:"id"`
	LastLogin string  `json:"last_login"`
	UserId    uint64  `json:"user_id"`
	Password  *string `json:"password,omitempty"`
}

type SessionRequest struct {
	Groupname string `json:"groupname"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

type Resource struct {
	Description       string `json:"description"`
	GroupID           int    `json:"group_id"`
	ID                int    `json:"id"`
	PrimaryFlag       bool   `json:"primary_flag"`
	RedactBookingText bool   `json:"redact_booking_text"`
	SequenceNum       int    `json:"sequence_num"`
	Symbol            string `json:"symbol"`
}

type User struct {
	Administrator         bool   `json:"administrator"`
	AllowAlterOthers      bool   `json:"allow_alter_others"`
	AllowEdits            bool   `json:"allow_edits"`
	BookingChangeEmails   bool   `json:"booking_change_emails"`
	ConfirmationEmailSent bool   `json:"confirmation_email_sent"`
	Email                 string `json:"email"`
	EmailConfirmed        bool   `json:"email_confirmed"`
	EmailMessages         bool   `json:"email_messages"`
	EmailPreferences      bool   `json:"email_preferences"`
	GroupID               int    `json:"group_id"`
	ID                    int    `json:"id"`
	LastLogin             string `json:"last_login"`
	MemberDetailsPrivate  bool   `json:"member_details_private"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
}

type Group struct {
	AviationRelated            bool   `json:"aviation_related"`
	DefaultTimeSlots           int    `json:"default_time_slots"`
	DefaultUserNameBookings    bool   `json:"default_user_name_bookings"`
	Description                string `json:"description"`
	FirstDayOfWeek             int    `json:"first_day_of_week"`
	Groupname                  string `json:"groupname"`
	ID                         int    `json:"id"`
	LastUpdateNum              string `json:"last_update_num"`
	Location                   string `json:"location"`
	OnHold                     bool   `json:"on_hold"`
	PrimaryResourceDescription string `json:"primary_resource_description"`
	RedactBookingText          string `json:"redact_booking_text"`
	ServiceUntil               string `json:"service_until"`
	Timezone                   string `json:"timezone"`
}

type SessionResonse struct {
	GroupID   uint64 `json:"group_id"`
	ID        string `json:"id"`
	LastLogin string `json:"last_login"`
	UserID    uint64 `json:"user_id"`
}

type BookedTimeSlot struct {
	BookingDate string `json:"booking_date"`
	GroupID     int    `json:"group_id"`
	ID          int    `json:"id"`
	TimeSlotID  int    `json:"time_slot_id"`
}

type Boooking struct {
	BlockUUID         interface{}   `json:"block_uuid"`
	BookedByUserID    int           `json:"booked_by_user_id"`
	BookedTimeSlotID  int           `json:"booked_time_slot_id"`
	BookedWhen        string        `json:"booked_when"`
	Description       string        `json:"description"`
	ID                int           `json:"id"`
	PrimaryBookingID  interface{}   `json:"primary_booking_id"`
	ResourceID        int           `json:"resource_id"`
	SecondaryBookings []interface{} `json:"secondary_bookings"`
}

type BookingRequest struct {
	ResourceID           int           `json:"resource_id"`
	Description          string        `json:"description"`
	BookedTimeSlotID     int           `json:"booked_time_slot_id"`
	BookedByUserID       int           `json:"booked_by_user_id"`
	BookedWhen           string        `json:"booked_when"`
	SecondaryResourceIds []interface{} `json:"secondary_resource_ids"`
	Ical                 bool          `json:"ical"`
}

type BookingResponse struct {
	BookingDate string `json:"booking_date"`
	GroupID     int    `json:"group_id"`
	ID          int    `json:"id"`
	TimeSlotID  int    `json:"time_slot_id"`
}

type BookingTimeSlotRequest struct {
	TimeSlotID  int    `json:"time_slot_id"`
	BookingDate string `json:"booking_date"`
}

type BookingTimeSlotResponse struct {
	BookingDate string `json:"booking_date"`
	GroupID     int    `json:"group_id"`
	ID          int    `json:"id"`
	TimeSlotID  int    `json:"time_slot_id"`
}

type TimeSlot struct {
	DayOfWeek   int    `json:"day_of_week"`
	Description string `json:"description"`
	GroupID     int    `json:"group_id"`
	ID          int    `json:"id"`
	Prime       bool   `json:"prime"`
	SequenceNum int    `json:"sequence_num"`
}
