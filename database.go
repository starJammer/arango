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
	json *databaseResult
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

type databaseResult struct {
	Result struct {
		Name     string
		Id       string
		Path     string
		IsSystem bool
	}
	ArangoError
}

func (db *Database) Name() string {
	return db.json.Result.Name
}

func (db *Database) Path() string {
	return db.json.Result.Path
}

func (db *Database) Id() string {
	return db.json.Result.Id
}

func (db *Database) IsSystem() bool {
	return db.json.Result.IsSystem
}

//UseDatabase will switch databases as if
//you called db._useDatabase in arangosh
//No error if successful, otherwise an erro
func (db *Database) UseDatabase(databaseName string) (*Database, error) {
	//create a new connection instead of re-using the old
	//object because re-use will cause collections
	//that used the old object to break
	return ConnDb(db.serverUrl.String(), databaseName)
}

//Small internal type used while creating a database
type createDatabase struct {
	Name  string `json:"name"`
	Users []User `json:"users"`
}

type createDatabaseResult struct {
	Result bool
	ArangoError
}

//CreateDatabase creates a new database and is modeled after db._createDatabase
//users can be nil, or it can be a list of users you want created
func (db *Database) CreateDatabase(name string, options *DatabaseOptions, users []User) error {

	var result createDatabaseResult
	var e ArangoError

	response, err := db.session.Post(db.serverUrl.String()+"/database", &createDatabase{Name: name, Users: users}, &result, &e)

	if err != nil {
		return newError(err.Error())
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
		return newError(err.Error())
	}

	switch response.Status() {
	case 200:
		return nil
	default:
		return e
	}

}

//Shortcut method for CreateDocumentCollection
//that will use default options to create the document
//collection.
//Similar to db._create( 'collection-name' ) in arangosh
func (db *Database) CreateDocumentCollection(collectionName string) error {
	return db.CreateCollection(collectionName, DefaultCollectionOptions())
}

//Shortcut method for CreateCollection that will
//use default options to create the edge collection
func (db *Database) CreateEdgeCollection(collectionName string) error {
	options := DefaultCollectionOptions()
	options.Type = EDGE_COLLECTION
	return db.CreateCollection(collectionName, options)
}

//CreateCollection is the generic collection creating method. Use it for more control.
//It allows you finer control by using the CollectionCreationOptions
func (db *Database) CreateCollection(collectionName string, options CollectionCreationOptions) error {

	options.Name = collectionName
	var e ArangoError

	response, err := db.session.Post(db.serverUrl.String()+"/collection", options, nil, e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200, 201:
		return nil
	default:
		return e
	}
	return nil
}

//Collection gets a collection by name from the database.
//Similar to calling db._collection( 'name' ) in arangosh
//See /_api/collection/{collection-name} endpoint for more info.
//An error with a code 404 is returned if it doesn't exist.
//
//Note : This method should not trigger the collection to be loaded
func (db *Database) Collection(collectionName string) (*Collection, error) {
	var c *Collection = new(Collection)
	c.db = db
	c.json = new(collectionResult)

	var e ArangoError

	response, err := db.session.Get(db.serverUrl.String()+"/collection/"+collectionName,
		nil,
		c.json,
		e,
	)

	if err != nil {
		return nil, newError(err.Error())
	}

	switch response.Status() {
	case 200, 201:
		return c, nil
	default:
		return nil, e
	}

	return nil, nil
}
