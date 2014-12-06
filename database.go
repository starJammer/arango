package arango

import (
	"errors"
	na "github.com/jmcvetta/napping"
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
//if something happened during the switch
func (db *Database) UseDatabase(dbName string) error {

	db.serverUrl.Path = "/_db/" + dbName + "/_api"
	var eMsg interface{}
	response, err := db.session.Get(db.serverUrl.String()+"/database/current", nil, db.json, &eMsg)

	if err != nil {
		return err
	}

	switch response.Status() {
	case 200:
		//The HTTP response is fine but arango returned an error
		if db.json.IsError {
			return errors.New(db.json.ErrorMessage)
		}

	//The HTTP reponse itself is a 401
	case 401:
		return errors.New("401 Unauthorized: check user password.")
	}

	return nil
}
//Used while creating a database
type createDatabase struct {
	Name  string `json:"name"`
	Users []User `json:"users"`
}

type createDatabaseResult struct {
	Result       bool
    basicJsonResult
}

type basicJsonResult struct {
    IsError        bool `json:"error"`
	Code         int
	ErrorNum     int
	ErrorMessage string
}

func (b *basicJsonResult) Error() string {
    return b.ErrorMessage
}

//CreateDatabase creates a new database and is modeled after db._createDatabase
func (db *Database) CreateDatabase(name string, options *DatabaseOptions, users []User) error {

	var result *createDatabaseResult = new(createDatabaseResult)
	var eMsg interface{}

	response, err := db.session.Post(db.serverUrl.String()+"/database", &createDatabase{Name: name, Users: users}, result, &eMsg)

	if err != nil {
        return ArangoError{ IsError : true, ErrorMessage : err.Error() }
	}

	switch response.Status() {
	case 201:
		return nil
    default:
        return result
	}

	return nil
}

func (db *Database) DropDatabase( name string ) error {

    var result *ArangoError = new(ArangoError)
    var eMsg interface{}

	response, err := db.session.Delete(db.serverUrl.String()+"/database/" + name, result, &eMsg)

    if err != nil {
        return ArangoError{ IsError : true, ErrorMessage : err.Error() }
    }

	switch response.Status() {
	case 201:
		return nil
    default:
        return result
	}

    return nil
}
