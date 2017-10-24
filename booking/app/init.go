package app

import (
	"github.com/revel/examples/booking/app/models"
	rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v2"
	"time"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}
	revel.OnAppStart(func() {
		Dbm := rgorp.Db.Map
		setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
			for col, size := range colSizes {
				t.ColMap(col).MaxSize = size
			}
		}

		t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
		t.ColMap("Password").Transient = true
		setColumnSizes(t, map[string]int{
			"Username": 20,
			"Name":     100,
		})

		t = Dbm.AddTable(models.Hotel{}).SetKeys(true, "HotelId")
		setColumnSizes(t, map[string]int{
			"Name":    50,
			"Address": 100,
			"City":    40,
			"State":   6,
			"Zip":     6,
			"Country": 40,
		})

		t = Dbm.AddTable(models.Booking{}).SetKeys(true, "BookingId")
		t.ColMap("User").Transient = true
		t.ColMap("Hotel").Transient = true
		t.ColMap("CheckInDate").Transient = true
		t.ColMap("CheckOutDate").Transient = true
		setColumnSizes(t, map[string]int{
			"CardNumber": 16,
			"NameOnCard": 50,
		})

		rgorp.Db.TraceOn(revel.AppLog)
		Dbm.CreateTables()

		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("demo"), bcrypt.DefaultCost)
		demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
		if err := Dbm.Insert(demoUser); err != nil {
			panic(err)
		}

		hotels := []*models.Hotel{
			&models.Hotel{0, "Marriott Courtyard", "Tower Pl, Buckhead", "Atlanta", "GA", "30305", "USA", 120},
			&models.Hotel{0, "W Hotel", "Union Square, Manhattan", "New York", "NY", "10011", "USA", 450},
			&models.Hotel{0, "Hotel Rouge", "1315 16th St NW", "Washington", "DC", "20036", "USA", 250},
		}
		for _, hotel := range hotels {
			if err := Dbm.Insert(hotel); err != nil {
				panic(err)
			}
		}
		bookings := []*models.Booking{
			&models.Booking{0, demoUser.UserId, hotels[0].HotelId, time.Now().Format(models.SQL_DATE_FORMAT), time.Now().Format(models.SQL_DATE_FORMAT), "id1", "n1", 12, 2, false, 2, time.Now(), time.Now(), demoUser, hotels[0]},
			&models.Booking{0, demoUser.UserId, hotels[1].HotelId, time.Now().Format(models.SQL_DATE_FORMAT), time.Now().Format(models.SQL_DATE_FORMAT), "id2", "n2", 12, 2, false, 2, time.Now(), time.Now(), demoUser, hotels[1]},
			&models.Booking{0, demoUser.UserId, hotels[2].HotelId, time.Now().Format(models.SQL_DATE_FORMAT), time.Now().Format(models.SQL_DATE_FORMAT), "id3", "n3", 12, 2, false, 2, time.Now(), time.Now(), demoUser, hotels[2]},
		}
		for _, booking := range bookings {
			if err := Dbm.Insert(booking); err != nil {
				panic(err)
			}
		}
	}, 5)
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
