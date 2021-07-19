package bookings

import (
	cld "Assignment/Packages/calendar"
	"time"

	dstore "Assignment/Packages/datastore"
	"fmt"
)

var BData []dstore.Bookings

//User add booking to data table
func CarBooking(car string, userName string, sD string, eD string) (int, error) {
	var days float64
	sDate := cld.StringToTime(sD)
	eDate := cld.StringToTime(eD)
	days = eDate.Sub(sDate).Hours() / 24
	if days <= 0 {
		days = 1
	}

	var p float64
	if days >= 30 {
		p = 1200 / 30
	} else if days >= 7 && days < 30 {
		p = 330 / 7
	} else {
		p = 55
	}

	//update booking and vehicle data
	Booking := dstore.PushBookings()
	Booking.UserName = userName
	Booking.CarName = car
	Booking.StartDate = sD
	Booking.EndDate = eD
	Booking.Price = p
	Booking.DaysOfRenting = int(days)
	row, err := dstore.InsertNewBookings(Booking)
	fmt.Println(row)
	return row, err
}

//Show User Bookings data
func ShowAllUserBookings(UserName string) [][]string {
	fmt.Println("bookings.ShowAllUserBookings")

	UserBookings, tf := dstore.GetDBUserBookings(UserName)
	if !tf {
		fmt.Println("no booking data for the user.")
		return nil
	}
	var bookings [][]string
	for _, v := range UserBookings {
		var booking []string
		tp := v.Price * float64(v.DaysOfRenting)
		id := fmt.Sprintf("CRA %v-%03d", time.Now().Year(), v.BookingID)
		days := fmt.Sprintf("%v", v.DaysOfRenting)
		price := fmt.Sprintf("$%.2f", v.Price)
		totalPrice := fmt.Sprintf("$%.2f", tp)
		booking = append(booking, id, v.CarName, v.StartDate, v.EndDate, days, price, totalPrice)
		bookings = append(bookings, booking)
	}
	return bookings
}
