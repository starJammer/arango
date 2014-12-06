package arango

import (
	"crypto/tls"
	"fmt"
	na "github.com/jmcvetta/napping"
	"net/http"
	"net/url"
)

var (
	//AllowBadSslCerts is to be used for development to allow self signed certs
	//This will not affect connections that have already been made, only
	//new connections that are created will be affected by this.
	AllowBadSslCerts = false
)

//Conn returns a new database connection to an arango server.
//This function will connect to the _system database
//
//The url can be in the following forms
//http://[username[:password]@]host:port
//or
//https://...etc
func Conn(host string) (*Database, error) {

	//connect to this db by default with these creds
	db := "_system"
	user := ""
	password := ""

	return ConnDbUserPassword(host, db, user, password)
}

//ConnDb returns a new database connection to an arango server.
//This function will connect to the given database using the
//default root user with a blank password.
func ConnDb(host, db string) (*Database, error) {

	user := ""
	password := ""

	return ConnDbUserPassword(host, db, user, password)
}

//ConnDbUserPassword returns a new database to an arango server.
//This function will connect to the given database using the
//given user name and password. This method is here for
//no reason really other than to separate the database
//host from the user name. So you can make host=http://localhost:1324
//and specify the user and password separately. Otherwise, you can
//just use ConnDb and specify the user info in the host string
func ConnDbUserPassword(host, dbName, user, password string) (*Database, error) {

	if dbName == "" {
		return nil, ArangoError{IsError: true, ErrorMessage: "A blank database was specified but that is not allowed."}
	}

	parsedUrl, err := url.Parse(host)

	if err != nil {
		return nil, err
	}

	if user != "" {
		parsedUrl.User = url.UserPassword(user, password)
	}

	switch parsedUrl.Scheme {
	case "http", "https":
	case "unix":
		return nil, ArangoError{
			IsError:      true,
			ErrorMessage: fmt.Sprintf("The %s scheme is not supported yet.", parsedUrl.Scheme),
		}
	default:
		return nil, ArangoError{
			IsError:      true,
			ErrorMessage: fmt.Sprintf("The %s scheme is not supported yet.", parsedUrl.Scheme),
		}
	}

	parsedUrl.Path = "/_db/" + dbName + "/_api"

	var db = new(Database)
	db.json = new(connectToDatabaseResult)
	db.serverUrl = parsedUrl
	db.session = new(na.Session)

	//Create transport that will allow
	//bad ssl certs
	if AllowBadSslCerts {
		db.session.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	var eMsg interface{}
	response, err := db.session.Get(db.serverUrl.String()+"/database/current", nil, db.json, &eMsg)

	if err != nil {
		return nil, err
	}

	switch response.Status() {
	case 200:
		//The HTTP response is fine but arango returned an error
		if db.json.IsError {
			return nil, db.json
		}

	//The HTTP reponse itself is a 401
	case 401:
		return nil, ArangoError{IsError: true, ErrorMessage: "401 Unauthorized: check user password."}
	}

	return db, nil

}
