package models

import (
	"github.com/revel/revel"
)

type Hotel struct {
	HotelID          int
	Name, Address    string
	City, State, Zip string
	Country          string
	Price            int
}

func (hotel *Hotel) Validate(v *revel.Validation) {
	v.Check(hotel.Name,
		revel.Required{},
		revel.MaxSize{Max: 50},
	)

	v.MaxSize(hotel.Address, 100)

	v.Check(hotel.City,
		revel.Required{},
		revel.MaxSize{Max: 40},
	)

	v.Check(hotel.State,
		revel.Required{},
		revel.MaxSize{Max: 6},
		revel.MinSize{Min: 2},
	)

	v.Check(hotel.Zip,
		revel.Required{},
		revel.MaxSize{Max: 6},
		revel.MinSize{Min: 5},
	)

	v.Check(hotel.Country,
		revel.Required{},
		revel.MaxSize{Max: 40},
		revel.MinSize{Min: 2},
	)
}
