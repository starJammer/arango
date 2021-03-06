package arango

import (
	"crypto/tls"
	"fmt"
	na "github.com/jmcvetta/napping"
	"net"
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
func ConnDb(host, databaseName string) (*Database, error) {

	user := ""
	password := ""

	return ConnDbUserPassword(host, databaseName, user, password)
}

type unixDialer struct {
	net.Dialer
    path string
}

// overriding net.Dialer.Dial to force unix socket connection
func (d *unixDialer) Dial(network, address string) (net.Conn, error) {
	return d.Dialer.Dial("unix", d.path)
}

//ConnDbUserPassword returns a new database to an arango server.
//This function will connect to the given database using the
//given user name and password. This method is here for
//no reason really other than to separate the database
//host from the user name. So you can make host=http://localhost:1324
//and specify the user and password separately. Otherwise, you can
//just use ConnDb and specify the user info in the host string
func ConnDbUserPassword(host, databaseName, user, password string) (*Database, error) {

	if databaseName == "" {
		return nil, ArangoError{IsError: true, ErrorMessage: "A blank database was specified but that is not allowed."}
	}

	parsedUrl, err := url.Parse(host)

	if err != nil {
		return nil, err
	}

	if user != "" {
		parsedUrl.User = url.UserPassword(user, password)
	}

	var db = new(Database)
	db.json = new(databaseResult)
	db.serverUrl = parsedUrl
	db.originalUrl = &url.URL{
        Scheme: parsedUrl.Scheme,
        Opaque: parsedUrl.Opaque,
        User: parsedUrl.User,
        Host: parsedUrl.Host,
        Path: parsedUrl.Path,
        RawQuery: parsedUrl.RawQuery,
        Fragment: parsedUrl.Fragment,
    }
	db.session = new(na.Session)

	switch parsedUrl.Scheme {
	case "http", "https":

		//Create transport that will allow
		//bad ssl certs if we allow it
		db.session.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: AllowBadSslCerts,
				},
			},
		}
	case "unix":
		//Create transport that will allow
		//bad ssl certs if we allow it
		db.session.Client = &http.Client{
			Transport: &http.Transport{
                Dial: (&unixDialer{
                        net.Dialer{},
                        parsedUrl.Path,
                    }).Dial,
			},
		}
        parsedUrl.Scheme = "http"
        parsedUrl.Host = "localhost-socket"
	default:
		return nil, newError(fmt.Sprintf("The %s scheme is not supported yet.", parsedUrl.Scheme))
	}

    parsedUrl.Path = "/_db/" + databaseName + "/_api"

	var e ArangoError
	response, err := db.session.Get(db.serverUrl.String()+"/database/current", nil, db.json, &e)

	if err != nil {
		return nil, newError(err.Error())
	}

	switch response.Status() {
	case 200:
		return db, nil
	case 401:
		return nil, newError("401 Unauthorized: check user password.")
	default:
		return nil, e
	}

}
