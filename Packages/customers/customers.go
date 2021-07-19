package customers

import (
	"fmt"
	"strings"

	dstore "Assignment/Packages/datastore"
)

func ExistingCustomer(input string) (Customers dstore.User, ok bool) {
	fmt.Println("customers.ExistingCustomer")
	fmt.Println(input)
	if len(input) == 36 && strings.Count(input, "-") == 4 {
		fmt.Println("customers.ExistingCustomer - match with uuid")
		Customers, ok = dstore.CheckDBUserKey(input)
		return
	} else {
		fmt.Println("customers.ExistingCustomer - match with user name")
		Customers, ok = dstore.CheckDBUserName(input)
		fmt.Println(Customers, ok)
		return
	}
	return dstore.User{}, false
}

func PasswordMatched(username, password string) bool {
	fmt.Println("customers.PasswordMatched")

	ok := dstore.VerifyPassword(username, password)
	return ok
}

func AddCustomer(username, password, key string) error {
	fmt.Println("customers.AddCustomer")

	row, err := dstore.InsertNewUsertoRow(username, password, key)
	fmt.Println(row)
	return err
}

func ChangePassword(username, newpassword string) error {
	fmt.Println("customers.ChangePassword")

	err := dstore.UpdateUserPassword(username, newpassword)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
