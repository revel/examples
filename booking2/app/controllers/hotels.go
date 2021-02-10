package controllers

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/revel/examples/booking2/app/models"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Hotels struct {
	Application
}

func (c Hotels) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(Application.Index)
	}
	return nil
}

func (c Hotels) Index() revel.Result {
	c.Log.Info("Fetching index")
	var bookings []*models.Booking
	_, err := c.Txn.Select(&bookings,
		c.Db.SqlStatementBuilder.Select("*").
			From("Booking").Where("UserID = ?", c.connected().UserID))
	if err != nil {
		panic(err)
	}

	return c.Render(bookings)
}

func (c Hotels) List(search string, size, page uint64) revel.Result {
	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	search = strings.TrimSpace(search)

	var hotels []*models.Hotel
	builder := c.Db.SqlStatementBuilder.Select("*").From("Hotel").Offset((page - 1) * size).Limit(size)
	if search != "" {
		search = "%" + strings.ToLower(search) + "%"
		builder = builder.Where(squirrel.Or{
			squirrel.Expr("lower(Name) like ?", search),
			squirrel.Expr("lower(City) like ?", search),
		})
	}
	if _, err := c.Txn.Select(&hotels, builder); err != nil {
		c.Log.Fatal("Unexpected error loading hotels", "error", err)
	}

	return c.Render(hotels, search, size, page, nextPage)
}

func (c Hotels) loadHotelByID(id int) *models.Hotel {
	h, err := c.Txn.Get(models.Hotel{}, id)
	if err != nil {
		panic(err)
	}

	if h == nil {
		return nil
	}

	return h.(*models.Hotel)
}

func (c Hotels) Show(id int) revel.Result {
	hotel := c.loadHotelByID(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}
	title := hotel.Name
	return c.Render(title, hotel)
}

func (c Hotels) Settings() revel.Result {
	return c.Render()
}

func (c Hotels) SaveSettings(password, verifyPassword string) revel.Result {
	models.ValidatePassword(c.Validation, password)
	c.Validation.Required(verifyPassword).
		Message("Please verify your password")
	c.Validation.Required(verifyPassword == password).
		Message("Your password doesn't match")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(Hotels.Settings)
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := c.Txn.ExecUpdate(c.Db.SqlStatementBuilder.
		Update("User").Set("HashedPassword", bcryptPassword).
		Where("UserID=?", c.connected().UserID))
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Password updated")
	return c.Redirect(Hotels.Index)
}

func (c Hotels) ConfirmBooking(id int, booking models.Booking) revel.Result {
	hotel := c.loadHotelByID(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}

	title := fmt.Sprintf("Confirm %s booking", hotel.Name)
	booking.Hotel = hotel
	booking.User = c.connected()
	booking.Validate(c.Validation)

	if c.Validation.HasErrors() || c.Params.Get("revise") != "" {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Hotels.Book, id)
	}

	if c.Params.Get("confirm") != "" {
		err := c.Txn.Insert(&booking)
		if err != nil {
			panic(err)
		}
		c.Flash.Success("Thank you, %s, your confirmation number for %s is %d",
			booking.User.Name, hotel.Name, booking.BookingID)
		return c.Redirect(Hotels.Index, id)
	}

	return c.Render(title, hotel, booking)
}

func (c Hotels) CancelBooking(id int) revel.Result {
	_, err := c.Txn.Delete(&models.Booking{BookingID: id})
	if err != nil {
		panic(err)
	}
	c.Flash.Success(fmt.Sprintln("Booking cancelled for confirmation number", id))
	return c.Redirect(Hotels.Index)
}

func (c Hotels) Book(id int) revel.Result {
	hotel := c.loadHotelByID(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}

	title := "Book " + hotel.Name
	return c.Render(title, hotel)
}
