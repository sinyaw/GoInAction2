package vehicle

import (
	dstore "Assignment/Packages/datastore"
)

//Get Name from TypeVehicle' Table on the database
func GetCarNames() []string {
	names := dstore.GetDBAllVehicalName()
	return names
}
