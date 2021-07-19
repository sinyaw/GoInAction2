package internet

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	bks "Assignment/Packages/bookings"
	cld "Assignment/Packages/calendar"
	ctms "Assignment/Packages/customers"
	vhcs "Assignment/Packages/vehicle"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	UserName string
	// Password string
	Key string
}

var tpl *template.Template
var emailRegex = regexp.MustCompile("^[\\w!#$%&'*+/=?`{|}~^-]+(?:\\.[\\w!#$%&'*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$") // regular expression
var mapUsers = map[string]user{}
var mapSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	log.SetFlags(log.Lshortfile)
}

func Htmlmain() {
	fmt.Println("Htmlmain")

	//Main Pages
	http.HandleFunc("/", mainMenu)
	http.HandleFunc("/menu", mainMenu)
	http.HandleFunc("/newuser", newUser)
	http.HandleFunc("/logout", logOut)
	http.HandleFunc("/login", logIn)
	//Sub Pages
	http.HandleFunc("/changepassword", userUpdatedPassword)
	http.HandleFunc("/bookcar", bookCar)
	http.HandleFunc("/viewbooking", viewBookingCarPage)
	//Load Files
	http.Handle("/Pictures/", http.StripPrefix("/Pictures", http.FileServer(http.Dir("Pictures"))))
	http.Handle("/Stuff/", http.StripPrefix("/Stuff", http.FileServer(http.Dir("Stuff"))))
	//run server
	log.Println(http.ListenAndServe(":5221", nil))

}

//Web Main Pages func from below---------------------------------------------------------------------------//
//new user
func newUser(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Htmlmain.newUser")

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "New User Registration", UserName: ""}
	tpl.ExecuteTemplate(res, "NewUser.gohtml", Data)

	// var currentUser user
	if req.Method == http.MethodPost {
		//---save user's information in the map ---
		username := req.FormValue("username")
		password := req.FormValue("password")
		confirmpassword := req.FormValue("confirmpassword")
		_, nameFound := ctms.ExistingCustomer(username)
		UFok := isEmailValid(username)
		if !nameFound && password == confirmpassword && UFok {
			// fmt.Println("Htmlmain.newUser - no same name in existing data")
			id := uuid.NewV4().String()
			setCookie(res, req, id)
			ctms.AddCustomer(username, password, id)
			mapSessions[id] = username
			mapUsers[username] = user{username, id}
			Data.MsgToUser = "New User Registration Done! You may process to log in."
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else if !UFok {
			// fmt.Println("Htmlmain.newUser - email format not correct")
			Data.MsgToUser = "Please enter correct email!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else if nameFound {
			// fmt.Println("Htmlmain.newUser - name in existing data")
			fmt.Scanf(Data.MsgToUser, "Please use other user name! '%v' has been taken!", username)
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else if password != confirmpassword {
			// fmt.Println("Htmlmain.newUser - confirm password not match")
			Data.MsgToUser = "Confirm Password is not same!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}
	}

}

//log in
func logIn(res http.ResponseWriter, req *http.Request) {
	fmt.Println("HTMLmain.logIn")

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "Log In"}
	myCookie, err := req.Cookie("CRA")
	if err == nil {
		Data.UserName = mapSessions[myCookie.Value]
	} else {
		Data.UserName = ""
	}

	if req.Method == http.MethodPost {
		userName := req.FormValue("username")
		password := req.FormValue("password")

		if userName == "" {
			Data.MsgToUser = "No value found in User Name!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else {
			matchUser, mOk := ctms.ExistingCustomer(userName)
			if pOk := ctms.PasswordMatched(matchUser.UserName, password); mOk && pOk {
				currentUser := user{userName, matchUser.Key}
				mapUsers[userName] = currentUser
				mapSessions[currentUser.Key] = userName
				// fmt.Println("Htmlmain.logIn - currentUser", currentUser)
				setCookie(res, req, currentUser.Key)
				http.Redirect(res, req, "/menu", http.StatusSeeOther)
			} else {
				Data.MsgToUser = "•	The User Name or Password is incorrect!"
				defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			}
		}
	}
	tpl.ExecuteTemplate(res, "Login.gohtml", Data)
}

//Main Menu
func mainMenu(res http.ResponseWriter, req *http.Request) {
	log.Println("mainMenu")

	Data := struct {
		PageName string
		UserName string
	}{PageName: "Main Menu"}
	var currentUser user
	myCookie, err := req.Cookie("CRA")
	if err != nil {
		// fmt.Println("Htmlmain.mainMenu - Cookie cannot find")
		http.Redirect(res, req, "/login", http.StatusSeeOther)
	} else if err == nil {
		// fmt.Println("Htmlmain.mainMenu - Cookie found")
		id := myCookie.Value
		matchUser, ok := ctms.ExistingCustomer(id)
		if !ok {
			// fmt.Println("Htmlmain.MainMenu - cookie fount with no record match in data")
			http.Redirect(res, req, "/login", http.StatusSeeOther)
		} else {
			// fmt.Println("Htmlmain.MainMenu - cookie fount with matched record in data")
			currentUser = user{matchUser.UserName, matchUser.Key}
			mapSessions[myCookie.Value] = matchUser.UserName
			mapUsers[matchUser.UserName] = currentUser
			Data.UserName = currentUser.UserName

		}
	}

	tpl.ExecuteTemplate(res, "MainMenu.gohtml", Data)
}

//Log Out
func logOut(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Htmlmain.logOut")

	Data := struct {
		PageName string
		UserName string
	}{PageName: "Log Out", UserName: "bye-bye"}
	Cookie, err := req.Cookie("CRA")
	if err == nil {
		Cookie.MaxAge = -1
		delete(mapUsers, mapSessions[Cookie.Value])
		delete(mapSessions, Cookie.Value)
		http.SetCookie(res, Cookie)
		// fmt.Println("Cookie deleted")
	} else {
		// fmt.Println("No Cookie found and deleted")
		http.Redirect(res, req, "/logIn", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(res, "LogOut.gohtml", Data)
}

//Web Sub Pages func start from below---------------------------------------------------------------------------//
//new user
func userUpdatedPassword(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Htmlmain.userUpdatedPassword")

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "Updated Password"}

	myCookie, err := req.Cookie("CRA")
	if err != nil {
		http.Redirect(res, req, "/logIn", http.StatusSeeOther)
	} else {
		Data.UserName = mapSessions[myCookie.Value]
	}
	if req.Method == http.MethodPost {
		//get user name and current password
		username := req.FormValue("username")
		password := req.FormValue("oldpassword")
		//get user new password and confirm the new password
		newpassword := req.FormValue("newpassword")
		confirmpassword := req.FormValue("confirmpassword")
		_, mOk := ctms.ExistingCustomer(username)
		pOk := ctms.PasswordMatched(username, password)
		fmt.Println("here")
		fmt.Println(username, password, mOk, pOk)
		if !mOk || !pOk {
			Data.MsgToUser = "The user name and password is not match! "
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			// http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
		} else if newpassword != confirmpassword {
			Data.MsgToUser = "New password and confrim password is not same!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			// http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
		} else if username == "" || password == "" {
			Data.MsgToUser = "Please insert username and password for verification!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else {
			ctms.ChangePassword(username, newpassword)
			Data.MsgToUser = "Password is updated!"
			defer fmt.Fprintf(res, "<h4 class='Application'><a href='/menu'>Main Menu</a></h4><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}
	}
	tpl.ExecuteTemplate(res, "ChangePassword.gohtml", Data)
}

//Book A Car
func bookCar(res http.ResponseWriter, req *http.Request) {
	fmt.Println("bookCar")

	Data := struct {
		PageName   string
		UserName   string
		CarDisplay []string
		Today      string
	}{PageName: "Rent A Car", Today: cld.TimeToString(time.Now())}

	cookie, err := req.Cookie("CRA")
	Data.UserName = mapSessions[cookie.Value]
	Data.CarDisplay = vhcs.GetCarNames()
	// fmt.Println("Htmlmain.bookCar - check user name: ",mapSessions[cookie.Value])
	if _, ok := mapSessions[cookie.Value]; err != nil || !ok {
		// tpl.Execute(res, "MainMenu.gohtml")
		http.Redirect(res, req, "/menu", http.StatusSeeOther)
		return
	} else {
		if req.Method == http.MethodPost {
			car := req.FormValue("car")
			startdate := req.FormValue("startdate")
			enddate := req.FormValue("enddate")
			// fmt.Println(car, startdate, enddate)
			if car != "" && startdate != "" && enddate != "" {
				bks.CarBooking(car, mapSessions[cookie.Value], startdate, enddate)
				http.Redirect(res, req, "/menu", http.StatusSeeOther)
				// bks.StoreData() // store data to excel
				return
			}
		}
		tpl.ExecuteTemplate(res, "BookCar.gohtml", Data)
	}
}

//Current Booking
func viewBookingCarPage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Htmlmain.viewBookingCarPage")

	Data := struct {
		PageName string
		UserName string
		Bookings [][]string
	}{PageName: "Current Booking"}
	cookie, err := req.Cookie("CRA")

	if _, ok := mapSessions[cookie.Value]; err != nil || !ok {
		http.Redirect(res, req, "/menu", http.StatusSeeOther)
	} else {
		currentUser := mapUsers[mapSessions[cookie.Value]]
		Data.UserName = currentUser.UserName
		Data.Bookings = bks.ShowAllUserBookings(Data.UserName)
		// fmt.Println(Data.Bookings)
	}
	tpl.ExecuteTemplate(res, "ViewBooking.gohtml", Data)
}

//Web pull out func ---------------------------------------------------------------------------//
//set cookie on client computer
func setCookie(res http.ResponseWriter, req *http.Request, id string) error {
	//name of cookies = "cookie" for 1hr & "CRA" for 2yrs
	co, _ := req.Cookie("CRA")
	co = &http.Cookie{
		Name:     "CRA",
		Value:    id,
		HttpOnly: false,
		Expires:  time.Now().AddDate(2, 0, 0),
	}
	http.SetCookie(res, co)
	// fmt.Println("Htmlmain.setCookie - done with set id = ", id)
	return nil
}

func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
