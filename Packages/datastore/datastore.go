package datastore

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
	log.SetFlags(log.Lshortfile)
}

//-------Function for DataBase---------------------------------------------------------------------------//
//start SQL Server
func RunMySQL() {
	fmt.Println("datastore.RunMySQL")

	var err error
	db, err = sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/mystoredb")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened.")
	}
}

//close SQL Server
func CloseMySQL() {
	db.Close()
}

//-------Function for TypeVehicle' Table---------------------------------------------------------------------------//
//get all the name column feild
func GetDBAllVehicalName() []string {
	fmt.Println("datastore.GetDBAllVehicalName")

	var names []string
	results, err := db.Query("SELECT Name FROM mystoredb.typevehicle")
	if err != nil {
		log.Println(err)
		return nil
	}
	for results.Next() {
		var name string
		err = results.Scan(&name)
		if err != nil {
			log.Println(err)
		} else {
			names = append(names, name)
		}
	}
	return names
}

//-------Function for Users' Table---------------------------------------------------------------------------//
type User struct {
	UserName string
	Key      string
}

// store new user to DB
func InsertNewUsertoRow(username, password, key string) (int, error) {
	fmt.Println("datastore.InsertNewUsertoRow")

	pw, err := hashPassword(password)
	if err != nil {
		return -1, err
	}

	results, err := db.Exec("INSERT INTO mystoredb.users VALUES (?, ?, ?)", username, pw, key)
	if err != nil {
		return -1, err
	} else {
		rows, _ := results.RowsAffected()
		return int(rows), nil
	}
}

// update user password
func UpdateUserPassword(username, newpassword string) error {
	fmt.Println("datastore.UpdateUserPassword")

	npw, err := hashPassword(newpassword)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE mystoredb.users SET Password = ? WHERE UserName = ?", npw, username)
	if err != nil {
		return err
	}
	return nil
}

// check existing name from Users Table
func CheckDBUserName(username string) (User, bool) {
	fmt.Println("datastore.CheckDBUserName")

	var user User
	var Password string
	results, err := db.Query("SELECT * FROM mystoredb.users WHERE UserName = ?", username)
	if err != nil {
		fmt.Println(err)
		return user, false
	}
	for results.Next() {
		err = results.Scan(&user.UserName, &Password, &user.Key)
		if err != nil {
			fmt.Println(err)
			return user, false
		} else {
			return user, true
		}
	}
	return user, false
}

// check existing key from Users Table
func CheckDBUserKey(Key string) (User, bool) {
	fmt.Println("datastore.CheckDBUserKey")

	var user User
	var Password string
	results, err := db.Query("Select * FROM mystoredb.users WHERE UUIDKey = ?", Key)
	if err != nil {
		fmt.Println(err)
		return user, false
	}
	for results.Next() {
		err = results.Scan(&user.UserName, &Password, &user.Key)
		if err != nil {
			fmt.Println(err)
			return user, false
		} else {
			return user, true
		}
	}
	return user, false
}

// hash the given password using bcrypt()
func hashPassword(password string) (string, error) {
	fmt.Println("datastore.hashPassword")

	if hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost); err != nil {
		log.Println(err)
		return "", err
	} else {
		return string(hash), nil
	}
}

// verify password from user input
func VerifyPassword(username, password string) bool {
	fmt.Println("datastore.VerifyPassword")

	var user User
	var hashedPassword string
	results, _ := db.Query("Select * FROM mystoredb.users WHERE UserName = ?", username)
	for results.Next() {
		results.Scan(&user.UserName, &hashedPassword, &user.Key)
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//-------Function for BookingsData' Table---------------------------------------------------------------------------//
type Bookings struct {
	BookingID     int
	UserName      string
	CarName       string
	StartDate     string
	EndDate       string
	Price         float64
	DaysOfRenting int
}

func PushBookings() Bookings {
	var bks Bookings
	return bks
}

// store new booking to DB
func InsertNewBookings(bk Bookings) (int, error) {
	fmt.Println("datastore.InsertNewUsertoRow")
	// fmt.Println(bk.UserName, bk.CarName, bk.StartDate, bk.EndDate, bk.Price, bk.DaysOfRenting)

	tableInput, err := db.Prepare("INSERT INTO mystoredb.BookingsData (UserName, CarName, StartDate, EndDate, Price, DaysOfRenting) Values(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	rowInput, err := tableInput.Exec(bk.UserName, bk.CarName, bk.StartDate, bk.EndDate, bk.Price, bk.DaysOfRenting)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := rowInput.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return int(rows), nil
}

// check bookings for the user from bookingsdata's Table
func GetDBUserBookings(UserName string) ([]Bookings, bool) {
	fmt.Println("datastore.GetDBUserBookings")

	var UserBookings []Bookings
	results, err := db.Query("SELECT * FROM mystoredb.BookingsData WHERE UserName = ?", UserName)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	for results.Next() {
		var bk Bookings
		err = results.Scan(&bk.UserName, &bk.CarName, &bk.StartDate, &bk.EndDate, &bk.Price, &bk.DaysOfRenting, &bk.BookingID)
		if err != nil {
			fmt.Println(err)
		}
		UserBookings = append(UserBookings, bk)
	}

	// fmt.Println(UserBookings)

	return UserBookings, true
}
