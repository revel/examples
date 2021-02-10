package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/revel/examples/booking/app/models"
	rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"github.com/revel/revel/logger"
	"github.com/valyala/fasthttp"
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
	logger.LogFunctionMap["stdoutjson"] =
		func(c *logger.CompositeMultiHandler, options *logger.LogOptions) {
			// Set the json formatter to os.Stdout, replace any existing handlers for the level specified
			c.SetJson(os.Stdout, options)
		}
	revel.AddInitEventHandler(func(event revel.Event, i interface{}) revel.EventResponse {
		if event == revel.ENGINE_BEFORE_INITIALIZED {
			if revel.RunMode == "dev-fast" {
				revel.AddHTTPMux("/this/is/a/test", fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
					fmt.Fprintln(ctx, "Hi there, it worked", string(ctx.Path()))
					ctx.SetStatusCode(200)
				}))
				revel.AddHTTPMux("/this/is/", fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
					fmt.Fprintln(ctx, "Hi there, shorter prefix", string(ctx.Path()))
					ctx.SetStatusCode(200)
				}))
			} else {
				revel.AddHTTPMux("/this/is/a/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "Hi there, it worked", r.URL.Path)
					w.WriteHeader(200)
				}))
				revel.AddHTTPMux("/this/is/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "Hi there, shorter prefix", r.URL.Path)
					w.WriteHeader(200)
				}))
			}
		}

		return 0
	})

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
		if err := Dbm.CreateTables(); err != nil {
			panic(err)
		}

		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("demo"), bcrypt.DefaultCost)
		demoUser := &models.User{
			Name:           "Demo User",
			Username:       "demo",
			Password:       "demo",
			HashedPassword: bcryptPassword,
		}

		if err := Dbm.Insert(demoUser); err != nil {
			panic(err)
		}

		count, _ := rgorp.Db.SelectInt(rgorp.Db.SqlStatementBuilder.Select("count(*)").From("User"))
		if count > 1 {
			revel.AppLog.Panic("Unexpected multiple users", "count", count)
		}

		hotels := []*models.Hotel{
			{
				Name:    "Marriott Courtyard",
				Address: "Tower Pl, Buckhead",
				City:    "Atlanta",
				State:   "GA",
				Zip:     "30305",
				Country: "USA",
				Price:   120,
			},
			{
				Name:    "W Hotel",
				Address: "Union Square, Manhattan",
				City:    "New York",
				State:   "NY",
				Zip:     "10011",
				Country: "USA",
				Price:   450,
			},
			{
				Name:    "Hotel Rouge",
				Address: "1315 16th St NW",
				City:    "Washington",
				State:   "DC",
				Zip:     "20036",
				Country: "USA",
				Price:   250,
			},
		}
		for _, hotel := range hotels {
			if err := Dbm.Insert(hotel); err != nil {
				panic(err)
			}
		}
		bookings := []*models.Booking{
			{
				UserID:       demoUser.UserID,
				HotelID:      hotels[0].HotelID,
				CheckInStr:   time.Now().Format(models.SQLDateFormat),
				CheckOutStr:  time.Now().Format(models.SQLDateFormat),
				CardNumber:   "id1",
				NameOnCard:   "n1",
				CardExpMonth: 12,
				CardExpYear:  2,
				Smoking:      false,
				Beds:         2,
				CheckInDate:  time.Now(),
				CheckOutDate: time.Now(),
				User:         demoUser,
				Hotel:        hotels[0],
			},
			{
				UserID:       demoUser.UserID,
				HotelID:      hotels[1].HotelID,
				CheckInStr:   time.Now().Format(models.SQLDateFormat),
				CheckOutStr:  time.Now().Format(models.SQLDateFormat),
				CardNumber:   "id2",
				NameOnCard:   "n2",
				CardExpMonth: 12,
				CardExpYear:  2,
				Smoking:      false,
				Beds:         2,
				CheckInDate:  time.Now(),
				CheckOutDate: time.Now(),
				User:         demoUser,
				Hotel:        hotels[1],
			},
			{
				UserID:       demoUser.UserID,
				HotelID:      hotels[2].HotelID,
				CheckInStr:   time.Now().Format(models.SQLDateFormat),
				CheckOutStr:  time.Now().Format(models.SQLDateFormat),
				CardNumber:   "id3",
				NameOnCard:   "n3",
				CardExpMonth: 12,
				CardExpYear:  2,
				Smoking:      false,
				Beds:         2,
				CheckInDate:  time.Now(),
				CheckOutDate: time.Now(),
				User:         demoUser,
				Hotel:        hotels[2],
			},
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
