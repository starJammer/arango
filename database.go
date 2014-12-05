package arango

import (
    na "github.com/jmcvetta/napping"
    "net/url"
    "errors"
)

//Database is an arango database connection.
//It is connected to one database at a time.
//You can switch databases with the UseDatabase
//Method.
//Do NOT instantiate this yourself. Use one
//of the Conn/ConnDb/ConnDbUserPassword 
//methods instead.
type Database struct{
    json *databaseJsonResult
    //holds addresses in form http://[username[:pass]@]localhost:8529
    serverUrl *url.URL
    session *na.Session
}

type databaseJsonResult struct{
    Result struct{
        Name string
        Id string
        Path string
        IsSystem bool
    }
    Error bool
    Code int
    ErrorNum int
    ErrorMessage string
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
func (db *Database) UseDatabase( dbName string ) error {

    db.serverUrl.Path = "/_db/" + dbName + "/_api"
	var eMsg interface{}
	response, err := db.session.Get(db.serverUrl.String()+"/database/current", nil, db.json, &eMsg)

    if err != nil {
        return err
    }

	switch response.Status() {
	case 200:
		//The HTTP response is fine but arango returned an error
		if db.json.Error {
			return errors.New(db.json.ErrorMessage)
		}

	//The HTTP reponse itself is a 401
	case 401:
		return errors.New("401 Unauthorized: check user password.")
	}

    return nil
}
