package arango

import (
	na "github.com/jmcvetta/napping"
    //"log"
	"net/url"
)

//Database is an arango database connection.
//It is connected to one database at a time.
//You can switch databases with the UseDatabase
//Method.
//Do NOT instantiate this yourself. Use one
//of the Conn/ConnDb/ConnDbUserPassword
//methods instead.
type Database struct {
	json *connectToDatabaseResult
	//holds addresses in form http://[username[:pass]@]localhost:8529
	serverUrl *url.URL
	session   *na.Session
}

//DatabaseOptions currently has nothing in it but is left as a
//placeholder for future use when arangodb makes use of it
type DatabaseOptions struct{}

type User struct {
	Username string      `json:"username"`
	Passwd   string      `json:"passwd"`
	Active   bool        `json:"active"`
	Extra    interface{} `json:"extra"`
}

type connectToDatabaseResult struct {
	Result struct {
		Name     string
		Id       string
		Path     string
		IsSystem bool
	}
	ArangoError
}

func (d *Database) Name() string {
	return d.json.Result.Name
}

func (d *Database) Path() string {
	return d.json.Result.Path
}

func (d *Database) Id() string {
	return d.json.Result.Id
}

func (d *Database) IsSystem() bool {
	return d.json.Result.IsSystem
}

//UseDatabase will switch databases as if
//you called db._useDatabase in arangosh
//No error if successful, otherwise an erro
func (db *Database) UseDatabase(dbName string) (*Database, error) {
	//create a new connection instead of re-using the old
	//object because re-use will cause collections
	//that used the old object to break
	return ConnDb(db.serverUrl.String(), dbName)
}

//Used while creating a database
type createDatabase struct {
	Name  string `json:"name"`
	Users []User `json:"users"`
}

type createDatabaseResult struct {
	Result bool
	ArangoError
}

//CreateDatabase creates a new database and is modeled after db._createDatabase
func (db *Database) CreateDatabase(name string, options *DatabaseOptions, users []User) error {

	var result createDatabaseResult
    var e ArangoError

	response, err := db.session.Post(db.serverUrl.String()+"/database", &createDatabase{Name: name, Users: users}, &result, &e)

	if err != nil {
		return newError( err.Error() )
	}

	switch response.Status() {
	case 201:
		return nil
	default:
		return e
	}

	return nil
}

type dropDatabaseResult struct {
    Result bool
    ArangoError
}

func (db *Database) DropDatabase(name string) error {

    var result dropDatabaseResult
	var e ArangoError

	response, err := db.session.Delete(db.serverUrl.String()+"/database/"+name, &result, &e)

	if err != nil {
		return newError( err.Error() )
	}

	switch response.Status() {
	case 200:
		return nil
	default:
		return e
	}

}
