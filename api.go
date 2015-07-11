package arango

import (
	gr "github.com/starJammer/grestclient"
)

const (
	Databasepath = "/_db/%s"

	AqlfunctionEndPoint = "/_api/aqlfunction"
	BatchEndPoint       = "/_api/batch"
	DatabaseEndPoint    = "/_api/database"
	CollectionEndPoint  = "/_api/collection"
	CursorEndPoint      = "/_api/cursor"
	DocumentEndPoint    = "/_api/document"
	EdgeEndPoint        = "/_api/edge"
	EdgesEndPoint       = "/_api/edges"
	EndpointEndPoint    = "/_api/endpoint"
	ExplainEndPoint     = "/_api/explain"
	ExportEndPoint      = "/_api/export"
	GraphEndPoint       = "/_api/graph"
	ImportEndPoint      = "/_api/import"
	IndexEndPoint       = "/_api/index"
	JobEndPoint         = "/_api/job"
	LogEndPoint         = "/_api/log"
	QueryEndPoint       = "/_api/query"
	ReplicationEndPoint = "/_api/replication"
	SimpleEndPoint      = "/_api/simple"
	StructureEndPoint   = "/_api/structure"
	SystemEndPoint      = "/_api/system"
	TasksEndPoint       = "/_api/tasks"
	TransactionEndPoint = "/_api/transaction"
	TraversalEndPoint   = "/_api/traversal"
	UserEndPoint        = "/_api/user"
	VersionEndPoint     = "/_api/version"
	WalEndPoint         = "/_api/wal"
)

//Version information about arango. Should I make it a simple struct?
type Version interface {
	Server() string
	Version() string
	Details() map[string]string
}

//ArangoError represents an arango error.
//An implementation can return an ArangoError
//to provide more detailed information about
//why a call failed.
type ArangoError interface {
	IsError() bool
	Code() int
	ErrorNum() int
	ErrorMessage() string
	//Error is the error interface in go
	Error() string

	//Reserved for some operations that will return
	//the id, rev, or key of the document.
	//For example, /_api/document/{doc-handle} when it
	//return a 412 error
	Id() string
	Rev() string
	Key() string
}

//HasGrestClient represents something that has a grestclient. You can
//then fetch the grest client and configure it to your liking.
//The arango rest client implementation provided here uses
//a grestclient to do the HTTP rest calls.
type HasGrestClient interface {
	GetGrestClient() gr.Client
}

//Connection represents a RESTFUL gateway to an arangodb server
type Connection interface {
	//Version -> GET on /_api/version
	Version(bool) (Version, error)

	//Database returns a RESTFUL gateway to the database
	//endpoint. The url used would be the url for the connection
	//plus the adequate path for this database.
	//Ex. http://localhost:8529/_db/{name}/_api/database where
	//{name} is the passed in database name.
	Database(name string) Database
}

type Database interface {
	//GetConnection returns connection associated with this database.
	//It should be non-nil
	Connection() Connection

	//Name returns the name of the database that this endpoint accesses.
	//In other words, what was Connection.Database called with?
	Name() string

	//Collection gets the collection endpoint for the given
	//name
	Collection(name string) Collection

	//Get -> GET on /_api/database
	Get() ([]string, error)

	//GetUser -> GET on /_api/database/user
	GetUser() ([]string, error)

	//GetCurrent -> GET on /_api/database/current
	GetCurrent() (CurrentResult, error)

	//Post -> POST on /_api/database
	Post() error

	//Delete -> DELETE on /_api/database/{name}
	Delete(name string) error

	//GetCollections -> GET on  /_api/collection
	GetCollections()
	//PostCollection -> POST on /_api/collection
	PostCollection()
}

type CurrentResult struct {
	Name     string `json:name`
	Id       string `json:id`
	Path     string `json:path`
	IsSystem bool   `json:isSystem`
}

type Collection interface {
	//Name returns the name of the collection this endpoint is for
	Name() string

	//Database gets the related database endpoint for the collection
	Database() Database

	//Get -> GET on /_api/collection/{name}
	Get() error

	//GetProperties -> GET on /_api/collection/{name}/properties
	GetProperties() error
}

type Document interface {
}
