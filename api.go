package arango

import (
	gr "github.com/starJammer/grestclient"
)

type CollectionType int

//Collection types
const (
	DOCUMENT_COLLECTION CollectionType = 2
	EDGE_COLLECTION     CollectionType = 3
)

type CollectionStatus int

//Collection statuses
const (
	NEW_BORN_STATUS       CollectionStatus = 1
	UNLOADED_STATUS       CollectionStatus = 2
	LOADED_STATUS         CollectionStatus = 3
	BEING_UNLOADED_STATUS CollectionStatus = 4
	DELETED_STATUS        CollectionStatus = 5
)

const (
	Databasepath   = "/_db/%s"
	CollectionPath = "/_api/collection/%s"

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
	Post(*PostDatabaseOptions) error

	//Delete -> DELETE on /_api/database/{name}
	Delete(name string) error

	//GetCollections -> GET on  /_api/collection
	GetCollections(excludeSystemCollections bool) (CollectionDescriptors, error)

	//PostCollection -> POST on /_api/collection
	PostCollection(options *CollectionCreationOptions) error
}

//CollectionCreationOptions represent options when creating a new collection.
//Look at the documentation for the POST to /_api/collection
//for information on the defaults, optionals, and required
//attributes.
type CollectionCreationOptions struct {
	Name           string         `json:"name"`
	WaitForSync    bool           `json:"waitForSync,omitempty"`
	DoCompact      bool           `json:"doCompact"`
	JournalSize    int            `json:"journalSize,omitempty"`
	IsSystem       bool           `json:"isSystem,omitempty"`
	IsVolatile     bool           `json:"isVolatile,omitempty"`
	KeyOptions     *KeyOptions    `json:"keyOptions,omitempty"`
	Type           CollectionType `json:"type,omitempty"`
	NumberOfShards int            `json:"numberOfShards,omitempty"`
	ShardKeys      []string       `json:"shardKeys,omitempty"`
}

//KeyOptions stores information about how a collection's key is configured.
//It is used during collection creation to specify how the new collection's
//key should be setup.
type KeyOptions struct {
	Type          string `json:"type,omitempty"`
	AllowUserKeys bool   `json:"allowUserKeys"`
	Increment     int    `json:"increment"`
	Offset        int    `json:"offset"`
}

//DefaultCollectionOptions creates a default set of collection options
func DefaultCollectionOptions() *CollectionCreationOptions {
	return &CollectionCreationOptions{
		DoCompact: true,
		Type:      DOCUMENT_COLLECTION,
	}
}

type CollectionDescriptor interface {
	Id() string
	Name() string
	IsSystem() bool
	Status() CollectionStatus
	Type() CollectionType
}

type CollectionDescriptors []CollectionDescriptor

func (c CollectionDescriptors) Find(name string) CollectionDescriptor {
	for _, d := range c {
		if d.Name() == name {
			return d
		}
	}
	return nil
}

type CurrentResult interface {
	Id() string
	Name() string
	Path() string
	IsSystem() bool
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

	//Delete -> DELETE on /_api/collection/{name}
	Delete() error
}

type Document interface {
}

type PostDatabaseOptions struct {
	Name  string `json:"name"`
	Users []User `json:"users,omitempty"`
}

type User struct {
	Username string      `json:"username"`
	Passwd   string      `json:"passwd"`
	Active   bool        `json:"active"`
	Extra    interface{} `json:"extra"`
}
