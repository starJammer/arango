package arango

import (
	na "github.com/jmcvetta/napping"
	//"log"
	"fmt"
	"net/http"
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
//you called db._useDatabase in arangosh.
//No error if successful, otherwise an error.
//Under the hood it just makes another call to ConnDb
//Using the same credentials used for the original database
//and returns the results.
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

	endpoint := fmt.Sprintf("%s/database/%s", db.serverUrl.String(), name)

	response, err := db.session.Delete(endpoint, &result, &e)

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

//Shortcut method for CreateCollection
//that will use default options to create the document
//collection.
func (db *Database) CreateDocumentCollection(collectionName string) (*Collection, error) {
	return db.CreateCollection(collectionName, DefaultCollectionOptions())
}

//Shortcut method for CreateCollection that will
//use default options to create the edge collection
func (db *Database) CreateEdgeCollection(collectionName string) (*Collection, error) {
	options := DefaultCollectionOptions()
	options.Type = EDGE_COLLECTION
	return db.CreateCollection(collectionName, options)
}

//CreateCollection is the generic collection creating method. Use it for more control.
//It allows you finer control by using the CollectionCreationOptions
//An error is returned if there was an issue. Otherwise, it was a success
//and you can use db.Collection( collectionName ) to get the collection
//you just created and work with it.
func (db *Database) CreateCollection(collectionName string, options CollectionCreationOptions) (*Collection, error) {

	options.Name = collectionName
	var e ArangoError

	var c *Collection = new(Collection)
	c.db = db
	c.json = new(collectionResult)

	endpoint := fmt.Sprintf("%s/collection", db.serverUrl.String())

	response, err := db.session.Post(endpoint, options, c.json, &e)

	//fmt.Printf( "( %T, %+v )\n( %T, %+v )\n ( %T, %+v )\n",
	//response,response,
	//err, err,
	//e, e )

	if err != nil {
		return nil, newError(err.Error())
	}

	switch response.Status() {
	case 200, 201:
		return c, nil
	default:
		return nil, e
	}
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

	endpoint := fmt.Sprintf("%s/collection/%s", db.serverUrl.String(), collectionName)
	response, err := db.session.Get(endpoint,
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

type dropCollectionResult struct {
	Id string
	ArangoError
}

//DropCollection drops the collection in the database by name.
func (db *Database) DropCollection(collectionName string) error {

	var result dropCollectionResult
	var e ArangoError

	endpoint := fmt.Sprintf("%s/collection/%s", db.serverUrl.String(), collectionName)
	response, err := db.session.Delete(endpoint, &result, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200:
		return nil
	default:
		return e
	}

	return nil
}

//Saves a document using the POST /_api/document endpoint.
//Look at arango api docs for more info.
func (db *Database) SaveDocumentWithOptions(document interface{}, options *SaveOptions) error {

	if options == nil {
		return newError("You must provide save options when using the database.SaveWithOptions method.")
	}

	if options.Collection == "" {
		return newError("You must provide a collection name in the options when using database.SaveWithOptions.")
	}

	var e ArangoError

	endpoint := fmt.Sprintf("%s/document?collection=%s&createCollection=%v&waitForSync=%v",
		db.serverUrl.String(),
		options.Collection,
		options.CreateCollection,
		options.WaitForSync,
	)

	response, err := db.session.Post(endpoint, document, document, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200, 201, 202:
		return nil
	default:
		return e
	}
}

//Document looks for a document in the database
func (db *Database) Document(documentHandle interface{}, document interface{}) error {
	return db.DocumentWithOptions(documentHandle, document, nil)
}

//DocumentWithOptions looks for a document in the database
func (db *Database) DocumentWithOptions(documentHandle interface{}, document interface{}, options *GetOptions) error {

	var id string
	switch dh := documentHandle.(type) {
	case string:
		id = dh
	case HasArangoId:
		id = dh.Id()
	default:
		return newError("The document handle you passed in is not valid.")
	}

	if id == "" {
		return newError("You must specify a documentHandle when fetching a document.")
	}

    if rev, ok := documentHandle.(HasArangoRev); ok {
        if options == nil {
            options = &GetOptions{}
        }

        options.IfMatch = rev.Rev()
    }

	if options != nil {
		if db.session.Header == nil {
			db.session.Header = &http.Header{}
			defer func() { db.session.Header = nil }()
		}

		if options.IfNoneMatch != "" {
			db.session.Header.Add("If-None-Match", options.IfNoneMatch)
			defer func() { db.session.Header.Del("If-None-Match") }()
		}
		if options.IfMatch != "" {
			db.session.Header.Add("If-Match", options.IfMatch)
			defer func() { db.session.Header.Del("If-Match") }()
		}
	}

	var e ArangoError

	endpoint := fmt.Sprintf("%s/document/%s", db.serverUrl.String(), id)

	response, err := db.session.Get(endpoint, nil, document, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200, 304:
		return nil
	default:
		return e
	}
}

func (db *Database) ReplaceDocumentWithOptions(documentHandle, document interface{}, options *ReplaceOptions) error {

	var id string
	switch dh := documentHandle.(type) {
	case string:
		id = dh
	case HasArangoId:
		id = dh.Id()
	default:
		return newError("The document handle you passed in is not valid.")
	}

	if id == "" {
		return newError("You must specify a documentHandle when replacing a document.")
	}

    if rev, ok := documentHandle.(HasArangoRev); ok {
        if options == nil {
            options = DefaultReplaceOptions()
        }

        options.IfMatch = rev.Rev()
    }

	var query string

	if options != nil {
		if db.session.Header == nil {
			db.session.Header = &http.Header{}
			defer func() { db.session.Header = nil }()
		}

		if options.IfMatch != "" {
			db.session.Header.Add("If-Match", options.IfMatch)
			defer func() { db.session.Header.Del("If-Match") }()
		}

		query += fmt.Sprintf("?waitForSync=%v", options.WaitForSync)

		if options.Rev != "" {
			query += fmt.Sprintf("&rev=%s", options.Rev)
		}

		if options.Policy != "" {
			query += fmt.Sprintf("&policy=%s", options.Policy)
		}
	}

	var e ArangoError

	endpoint := fmt.Sprintf("%s/document/%s%s", db.serverUrl.String(), id, query)

	response, err := db.session.Put(endpoint, document, document, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200, 201, 202:
		return nil
	default:
		return e
	}

	return nil
}

func (db *Database) UpdateDocumentWithOptions(documentHandle, document interface{}, options *UpdateOptions) error {

	var id string
	switch dh := documentHandle.(type) {
	case string:
		id = dh
	case HasArangoId:
		id = dh.Id()
	default:
		return newError("The document handle you passed in is not valid.")
	}

	if id == "" {
		return newError("You must specify a documentHandle when updating a document.")
	}

    if rev, ok := documentHandle.(HasArangoRev); ok {
        if options == nil {
            options = DefaultUpdateOptions()
        }

        options.IfMatch = rev.Rev()
    }

	var query string

	if options != nil {
		if db.session.Header == nil {
			db.session.Header = &http.Header{}
			defer func() { db.session.Header = nil }()
		}

		if options.IfMatch != "" {
			db.session.Header.Add("If-Match", options.IfMatch)
			defer func() { db.session.Header.Del("If-Match") }()
		}

		query += fmt.Sprintf("?waitForSync=%v", options.WaitForSync)
		query += fmt.Sprintf("&keepNull=%v", options.KeepNull)
		query += fmt.Sprintf("&mergeArrays=%v", options.MergeArrays)

		if options.Rev != "" {
			query += fmt.Sprintf("&rev=%s", options.Rev)
		}

		if options.Policy != "" {
			query += fmt.Sprintf("&policy=%s", options.Policy)
		}
	}

	var e ArangoError

	endpoint := fmt.Sprintf("%s/document/%s%s", db.serverUrl.String(), id, query)

	response, err := db.session.Patch(endpoint, document, document, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 201, 202:
		return nil
	default:
		return e
	}

	return nil
}

func (db *Database) DeleteDocumentWithOptions(documentHandle interface{}, options *DeleteOptions) error {

	var id string
	switch dh := documentHandle.(type) {
	case string:
		id = dh
	case HasArangoId:
		id = dh.Id()
	default:
		return newError("The document handle you passed in is not valid.")
	}

	if id == "" {
		return newError("You must specify a documentHandle when deleting a document.")
	}

    if rev, ok := documentHandle.(HasArangoRev); ok {
        if options == nil {
            options = DefaultDeleteOptions()
        }

        options.IfMatch = rev.Rev()
    }

	var query string

	if options != nil {
		if db.session.Header == nil {
			db.session.Header = &http.Header{}
			defer func() { db.session.Header = nil }()
		}

		if options.IfMatch != "" {
			db.session.Header.Add("If-Match", options.IfMatch)
			defer func() { db.session.Header.Del("If-Match") }()
		}

		query += fmt.Sprintf("?waitForSync=%v", options.WaitForSync)

		if options.Rev != "" {
			query += fmt.Sprintf("&rev=%s", options.Rev)
		}

		if options.Policy != "" {
			query += fmt.Sprintf("&policy=%s", options.Policy)
		}
	}

	var e ArangoError

	endpoint := fmt.Sprintf("%s/document/%s%s", db.serverUrl.String(), id, query)

	response, err := db.session.Delete(endpoint, &struct{}{}, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200, 202:
		return nil
	default:
		return e
	}

	return nil
}
