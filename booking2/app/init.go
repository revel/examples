package app

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/revel/examples/booking2/app/models"
	rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

func init() { //nolint:gochecknoinits
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

		t := Dbm.AddTable(models.User{}).SetKeys(true, "UserID")
		t.ColMap("Password").Transient = true
		setColumnSizes(t, map[string]int{
			"Username": 20,
			"Name":     100,
		})

		t = Dbm.AddTable(models.Hotel{}).SetKeys(true, "HotelID")
		setColumnSizes(t, map[string]int{
			"Name":    50,
			"Address": 100,
			"City":    40,
			"State":   6,
			"Zip":     6,
			"Country": 40,
		})

		t = Dbm.AddTable(models.Booking{}).SetKeys(true, "BookingID")
		t.ColMap("User").Transient = true
		t.ColMap("Hotel").Transient = true
		t.ColMap("CheckInDate").Transient = true
		t.ColMap("CheckOutDate").Transient = true
		setColumnSizes(t, map[string]int{
			"CardNumber": 16,
			"NameOnCard": 50,
		})

		rgorp.Db.TraceOn(revel.AppLog)
		err := Dbm.CreateTables()
		if err != nil {
			revel.AppLog.Fatal("Failed to create tables", "error", err)
		}

		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("demo"), bcrypt.DefaultCost)
		demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
		if err := Dbm.Insert(demoUser); err != nil {
			panic(err)
		}
		count, _ := rgorp.Db.SelectInt(rgorp.Db.SqlStatementBuilder.Select("count(*)").From("User"))
		if count > 1 {
			revel.AppLog.Panic("Unexpected multiple users", "count", count)
		}

		hotels := []*models.Hotel{
			{0, "Marriott Courtyard", "Tower Pl, Buckhead", "Atlanta", "GA", "30305", "USA", 120},
			{0, "W Hotel", "Union Square, Manhattan", "New York", "NY", "10011", "USA", 450},
			{0, "Hotel Rouge", "1315 16th St NW", "Washington", "DC", "20036", "USA", 250},
		}
		for _, hotel := range hotels {
			if err := Dbm.Insert(hotel); err != nil {
				panic(err)
			}
		}
		bookings := []*models.Booking{
			{0, demoUser.UserID, hotels[0].HotelID, time.Now().Format(models.SQLDateFormat), time.Now().Format(models.SQLDateFormat), "id1", "n1", 12, 2, false, 2, time.Now(), time.Now(), demoUser, hotels[0]},
			{0, demoUser.UserID, hotels[1].HotelID, time.Now().Format(models.SQLDateFormat), time.Now().Format(models.SQLDateFormat), "id2", "n2", 12, 2, false, 2, time.Now(), time.Now(), demoUser, hotels[1]},
			{0, demoUser.UserID, hotels[2].HotelID, time.Now().Format(models.SQLDateFormat), time.Now().Format(models.SQLDateFormat), "id3", "n3", 12, 2, false, 2, time.Now(), time.Now(), demoUser, hotels[2]},
		}
		for _, booking := range bookings {
			if err := Dbm.Insert(booking); err != nil {
				panic(err)
			}
		}
	}, 5)
}

func HeaderFilter(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
