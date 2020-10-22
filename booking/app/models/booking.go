package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type Booking struct {
	BookingID    int
	UserID       int
	HotelID      int
	CheckInStr   string
	CheckOutStr  string
	CardNumber   string
	NameOnCard   string
	CardExpMonth int
	CardExpYear  int
	Smoking      bool
	Beds         int

	// Transient
	CheckInDate  time.Time
	CheckOutDate time.Time
	User         *User
	Hotel        *Hotel
}

// TODO: Make an interface for Validate() and then validation can pass in the
// key prefix ("booking.").
func (b Booking) Validate(v *revel.Validation) {
	v.Required(b.User)
	v.Required(b.Hotel)
	v.Required(b.CheckInDate)
	v.Required(b.CheckOutDate)

	v.Match(b.CardNumber, regexp.MustCompile(`\d{16}`)).
		Message("Credit card number must be numeric and 16 digits")

	v.Check(b.NameOnCard,
		revel.Required{},
		revel.MinSize{Min: 3},
		revel.MaxSize{Max: 70},
	)
}

func (b Booking) Total() int {
	return b.Hotel.Price * b.Nights()
}

func (b Booking) Nights() int {
	return int((b.CheckOutDate.Unix() - b.CheckInDate.Unix()) / 60 / 60 / 24)
}

const (
	DATE_FORMAT     = "Jan _2, 2006"
	SQL_DATE_FORMAT = "2006-01-02"
)

func (b Booking) Description() string {
	if b.Hotel == nil {
		return ""
	}

	return fmt.Sprintf("%s, %s to %s",
		b.Hotel.Name,
		b.CheckInDate.Format(DATE_FORMAT),
		b.CheckOutDate.Format(DATE_FORMAT))
}

func (b Booking) String() string {
	return fmt.Sprintf("Booking(%s,%s)", b.User, b.Hotel.Name)
}

// These hooks work around two things:
// - Gorp's lack of support for loading relations automatically.
// - Sqlite's lack of support for datetimes.

func (b *Booking) PreInsert(_ gorp.SqlExecutor) error {
	b.UserID = b.User.UserID
	b.HotelID = b.Hotel.HotelID
	b.CheckInStr = b.CheckInDate.Format(SQL_DATE_FORMAT)
	b.CheckOutStr = b.CheckOutDate.Format(SQL_DATE_FORMAT)
	return nil
}

func (b *Booking) PostGet(exe gorp.SqlExecutor) error {
	var (
		obj interface{}
		err error
	)

	obj, err = exe.Get(User{}, b.UserID)
	if err != nil {
		return fmt.Errorf("error loading a booking's user (%d): %s", b.UserID, err)
	}

	b.User = obj.(*User)

	obj, err = exe.Get(Hotel{}, b.HotelID)
	if err != nil {
		return fmt.Errorf("error loading a booking's hotel (%d): %s", b.HotelID, err)
	}
	b.Hotel = obj.(*Hotel)

	if b.CheckInDate, err = time.Parse(SQL_DATE_FORMAT, b.CheckInStr); err != nil {
		return fmt.Errorf("error parsing check in date '%s' %s", b.CheckInStr, err.Error())
	}
	if b.CheckOutDate, err = time.Parse(SQL_DATE_FORMAT, b.CheckOutStr); err != nil {
		return fmt.Errorf("error parsing check out date '%s' %s", b.CheckOutStr, err.Error())
	}
	return nil
}
