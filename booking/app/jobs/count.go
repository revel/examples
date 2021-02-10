package jobs

import (
	"fmt"

	"github.com/revel/examples/booking/app/models"
	"github.com/revel/modules/jobs/app/jobs"
	gorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

// Periodically count the bookings in the database.
type BookingCounter struct{}

func (c BookingCounter) Run() {
	bookings, err := gorp.Db.Map.Select(models.Booking{},
		`select * from Booking`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d bookings.\n", len(bookings))
}

func init() { //nolint:gochecknoinits
	revel.OnAppStart(func() {
		err := jobs.Schedule("@every 1m", BookingCounter{})
		if err != nil {
			panic(err)
		}
	})
}
