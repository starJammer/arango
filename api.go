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
	LOADING_STATUS        CollectionStatus = 6
)

type Policy string

const (
	ERROR_POLICY Policy = "error"
	LAST_POLICY  Policy = "last"
)

const (
	DatabasePrefix   = "/_db/%s"
	CollectionPrefix = "/_api/collection/%s"

	AqlfunctionPath = "/_api/aqlfunction"
	BatchPath       = "/_api/batch"
	DatabasePath    = "/_api/database"
	CollectionPath  = "/_api/collection"
	CursorPath      = "/_api/cursor"
	DocumentPath    = "/_api/document"
	EdgePath        = "/_api/edge"
	EdgesPath       = "/_api/edges"
	EndpointPath    = "/_api/endpoint"
	ExplainPath     = "/_api/explain"
	ExportPath      = "/_api/export"
	GraphPath       = "/_api/graph"
	ImportPath      = "/_api/import"
	IndexPath       = "/_api/index"
	JobPath         = "/_api/job"
	LogPath         = "/_api/log"
	QueryPath       = "/_api/query"
	ReplicationPath = "/_api/replication"
	SimplePath      = "/_api/simple"
	StructurePath   = "/_api/structure"
	SystemPath      = "/_api/system"
	TasksPath       = "/_api/tasks"
	TransactionPath = "/_api/transaction"
	TraversalPath   = "/_api/traversal"
	UserPath        = "/_api/user"
	VersionPath     = "/_api/version"
	WalPath         = "/_api/wal"
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
	//Error is for the error interface in go
	Error() string
	//usually true because the api doesn't return an error
	//if this is false
	IsError() bool
	//The http response code
	Code() int
	//arango error number
	ErrorNum() int
	//error description from arango server
	ErrorMessage() string

	//Reserved for some operations that will return
	//the id, rev, or key of the document in an error
	//json object.
	//This is primarily used when arango returns a 412
	//http respons code.
	//For example, GET /_api/document/{doc-handle} when it
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

	//Collection gets the collection endpoint.
	CollectionEndpoint() CollectionEndpoint

	//DocumentEndPoint gets the document endpoint
	DocumentEndpoint() DocumentEndpoint

	//EdgeEndPoint gets the document endpoint
	EdgeEndpoint() EdgeEndpoint

	//Get -> GET on /_api/database
	Get() ([]string, error)

	//GetUser -> GET on /_api/database/user
	GetUser() ([]string, error)

	//GetCurrent -> GET on /_api/database/current
	GetCurrent() (CurrentResult, error)

	//Post -> POST on /_api/database
	Post(name string, options *PostDatabaseOptions) error

	//Delete -> DELETE on /_api/database/{name}
	Delete(name string) error
}

//EdgeImplementation is an embeddable type that
//you can use to easily gain access to arango
//specific attributes for edges. These include
//the _id, _key, and _rev attributes from the
//DocumentImplementation as well as the _to and
//_from attributes for edges.
type EdgeImplementation struct {
	DocumentImplementation
	ArangoFrom string `json:"_from,omitempty"`
	ArangoTo   string `json:"_to,omitempty"`
}

func (e EdgeImplementation) From() string {
	return e.ArangoFrom
}

func (e EdgeImplementation) To() string {
	return e.ArangoTo
}

type EdgeEndpoint interface {
	//GetEdges -> GET on /_api/edge
	//If pased in returnType is "" then the default should be used.
	//Default is "path"
	GetEdges(collection string, options *GetEdgesOptions) ([]string, error)

	//PostEdge -> POST on /_api/edge
	//PostEdgeOptions are optional.
	//The same edge is populated with _id, _key, _rev attributes
	//if possible on a successful POST of the edge.
	PostEdge(edge interface{}, collection string, from string, to string, options *PostEdgeOptions) error

	//GetEdge -> GET on /_api/edge/{edge-handle}
	//edgeReceiver is where the edge will be json.Unmarshaled into.
	//GetEdgeOptions are optional and can be nil.
	//EdgeReceiver cannot be populated if you provide an
	//If-None-Match option and the server returns a 304 because the edge
	//revision matches If-None-Match. See arango docs for more info.
	//In this case, no error is returned and the edgeReceiver isn't
	//altered at all because the server didn't return any edge
	//attributes.
	//
	//EdgeReceiver IS NOT populated if you provide an
	//If-Match option and the server returns a 412 because the edge
	//revision does not match If-Match.  In this case, use the
	//error that is returned to get the latest revision number
	//by calling error.Rev()
	GetEdge(edgeHandle string, edgeReceiver interface{}, options *GetEdgeOptions) error

	//HeadEdge -> HEAD on /_api/edge/{edge-handle}
	//Returns the current revision of the edge.
	//If the edge doesn't exist the returned revision is blank.
	//In all other cases, the current revision is returned and
	//the error is nil.
	HeadEdge(edgeHandle string, options *HeadEdgeOptions) (revision string, err error)

	//PutEdge -> PUT on /_api/edge/{edge-handle}
	PutEdge(edgeHandle string, edge interface{}, options *PutEdgeOptions) error

	//PatchEdge -> PATCH on /_api/edge/{edge-handle}
	PatchEdge(edgeHandle string, edge interface{}, options *PatchEdgeOptions) error

	//DeleteEdge -> DELETE on /_api/edge/{edge-handle}
	DeleteEdge(edgeHandle string, options *DeleteEdgeOptions) error
}

type DocumentEndpoint interface {
	//GetDocuments -> GET on /_api/document
	//If pased in returnType is "" then the default should be used.
	//Default is "path"
	GetDocuments(collection string, options *GetDocumentsOptions) ([]string, error)

	//PostDocument -> POST on /_api/document
	//PostDocumentOptions are optional.
	//The same document is populated with _id, _key, _rev attributes
	//if possible on a successful POST of the document.
	PostDocument(document interface{}, collection string, options *PostDocumentOptions) error

	//GetDocument -> GET on /_api/document/{document-handle}
	//documentReceiver is where the document will be json.Unmarshaled into.
	//GetDocumentOptions are optional and can be nil.
	//DocumentReceiver cannot be populated if you provide an
	//If-None-Match option and the server returns a 304 because the document
	//revision matches If-None-Match. See arango docs for more info.
	//In this case, no error is returned and the documentReceiver isn't
	//altered at all because the server didn't return any document
	//attributes.
	//
	//DocumentReceiver IS NOT populated if you provide an
	//If-Match option and the server returns a 412 because the document
	//revision does not match If-Match.  In this case, use the
	//error that is returned to get the latest revision number
	//by calling error.Rev()
	GetDocument(documentHandle string, documentReceiver interface{}, options *GetDocumentOptions) error

	//HeadDocument -> HEAD on /_api/document/{document-handle}
	//Returns the current revision of the document.
	//If the document doesn't exist the returned revision is blank.
	//In all other cases, the current revision is returned and
	//the error is nil.
	HeadDocument(documentHandle string, options *HeadDocumentOptions) (revision string, err error)

	//PutDocument -> PUT on /_api/document/{document-handle}
	PutDocument(documentHandle string, document interface{}, options *PutDocumentOptions) error

	//PatchDocument -> PATCH on /_api/document/{document-handle}
	PatchDocument(documentHandle string, document interface{}, options *PatchDocumentOptions) error

	//DeleteDocument -> DELETE on /_api/document/{document-handle}
	DeleteDocument(documentHandle string, options *DeleteDocumentOptions) error
}

//DocumentImplementation is an embeddable type that
//you can use to easily gain access to arango specific
//data. It will help capture the _id, _key, and _rev
//attributes from responses made by  arangodb
type DocumentImplementation struct {
	ArangoId  string `json:"_id,omitempty"`
	ArangoRev string `json:"_rev,omitempty"`
	ArangoKey string `json:"_key,omitempty"`
}

func (d DocumentImplementation) Id() string {
	return d.ArangoId
}

func (d DocumentImplementation) Rev() string {
	return d.ArangoRev
}

func (d DocumentImplementation) Key() string {
	return d.ArangoKey
}

type PostDocumentOptions struct {
	CreateCollection bool //default is false
	WaitForSync      bool //default is false
}

type PostEdgeOptions PostDocumentOptions

type GetDocumentsOptions struct {
	ReturnType string
}
type GetEdgesOptions GetDocumentsOptions

type GetDocumentOptions struct {
	//Rev is used in the query
	Rev string
	//IfMatch is a header equivalent of IfMatch
	IfMatch string

	//IfNoneMatch is header
	IfNoneMatch string
}

type GetEdgeOptions GetDocumentOptions

type HeadDocumentOptions GetDocumentOptions

type HeadEdgeOptions HeadDocumentOptions

type PutDocumentOptions struct {
	WaitForSync bool
	Rev         string
	Policy      Policy
	IfMatch     string
}

type PutEdgeOptions PutDocumentOptions

//Use DefaultPatchDocumentOptions for options
//set to default arango values
type PatchDocumentOptions struct {
	KeepNull     bool
	MergeObjects bool

	WaitForSync bool
	Rev         string
	Policy      Policy
	IfMatch     string
}

type PatchEdgeOptions PatchDocumentOptions

func DefaultPatchDocumentOptions() *PatchDocumentOptions {
	return &PatchDocumentOptions{
		KeepNull:     true,
		MergeObjects: true,
		WaitForSync:  false,
	}
}

func DefaultPatchEdgeOptions() *PatchEdgeOptions {
	return &PatchEdgeOptions{
		KeepNull:     true,
		MergeObjects: true,
		WaitForSync:  false,
	}
}

type DeleteDocumentOptions struct {
	Rev         string
	Policy      Policy
	WaitForSync string
	IfMatch     string
}

type DeleteEdgeOptions DeleteDocumentOptions

type CollectionPropertyChange struct {
	WaitForSync bool `json:"waitForSync"`
	JournalSize int  `json:"journalSize,omitempty"`
}

//PostCollectionOptions represent options when creating a new collection.
//Look at the documentation for the POST to /_api/collection
//for information on the default, optional, and required
//attributes.
type PostCollectionOptions struct {
	Name               string              `json:"name"`
	WaitForSync        bool                `json:"waitForSync,omitempty"`
	DoCompact          bool                `json:"doCompact"`
	JournalSize        int                 `json:"journalSize,omitempty"`
	IsSystem           bool                `json:"isSystem,omitempty"`
	IsVolatile         bool                `json:"isVolatile,omitempty"`
	KeyCreationOptions *KeyCreationOptions `json:"keyOptions,omitempty"`
	Type               CollectionType      `json:"type,omitempty"`
	NumberOfShards     int                 `json:"numberOfShards,omitempty"`
	ShardKeys          []string            `json:"shardKeys,omitempty"`
}

//KeyOptions stores information about how a collection's key is configured.
//It is used during collection creation to specify how the new collection's
//key should be setup.
type KeyCreationOptions struct {
	Type          string `json:"type,omitempty"`
	AllowUserKeys bool   `json:"allowUserKeys"`
	Increment     int    `json:"increment"`
	Offset        int    `json:"offset"`
}

//DefaultPostCollectionOptions creates a default set of collection options
func DefaultPostCollectionOptions() *PostCollectionOptions {
	return &PostCollectionOptions{
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
	WaitForSync() bool
	DoCompact() bool
	JournalSize() int
	KeyOptions() KeyOptions
	IsVolatile() bool
	NumberOfShards() int
	ShardKeys() []string
	Count() int
	Figures() Figures
	Revision() string
	Checksum() int
}

type KeyOptions interface {
	Type() string
	AllowUserKeys() bool
	Increment() int
	Offset() int
}

type Figures interface {
	Alive() Alive
	Dead() Dead
	Datafiles() Datafiles
	Journals() Journals
	Compactors() Compactors
	Shapefiles() Shapefiles
	Shapes() Shapes
	Attributes() Attributes
	Indexes() Indexes
	MaxTick() string
	UncollectedLogfileEntries() int
}

type Alive interface {
	Count() int
	Size() int
}

type Dead interface {
	Count() int
	Size() int
	Deletion() int
}

type Datafiles interface {
	Count() int
	FileSize() int
}

type Journals interface {
	Count() int
	FileSize() int
}

type Compactors interface {
	Count() int
	FileSize() int
}

type Shapefiles interface {
	Count() int
	FileSize() int
}

type Shapes interface {
	Count() int
	Size() int
}

type Attributes interface {
	Count() int
	Size() int
}

type Indexes interface {
	Count() int
	Size() int
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

type CollectionEndpoint interface {

	//Database gets the related database endpoint
	//for this collection endpoint
	Database() Database

	//GetCollections -> GET on  /_api/collection
	GetCollections(excludeSystemCollections bool) (CollectionDescriptors, error)

	//PostCollection -> POST on /_api/collection
	PostCollection(name string, options *PostCollectionOptions) error

	//Get -> GET on /_api/collection/{name}
	Get(name string) (CollectionDescriptor, error)

	//GetProperties -> GET on /_api/collection/{name}/properties
	GetProperties(name string) (CollectionDescriptor, error)

	//GetCount -> GET on /_api/collection/{name}/count
	GetCount(name string) (CollectionDescriptor, error)

	//GetFigures -> GET on /_api/collection/{name}/figures
	GetFigures(name string) (CollectionDescriptor, error)

	//GetRevision -> GET on /_api/collection/{name}/revision
	GetRevision(name string) (CollectionDescriptor, error)

	//GetChecksum -> GET on /_api/collection/{name}/checksum
	GetChecksum(name string, withRevisions bool, withData bool) (CollectionDescriptor, error)

	//PutLoad -> PUT on /_api/collection/{name}/load
	PutLoad(name string, includeCount bool) (CollectionDescriptor, error)

	//PutUnload -> PUT on /_api/collection/{name}/unload
	PutUnload(name string) (CollectionDescriptor, error)

	//PutTruncate -> PUT on /_api/collection/{name}/truncate
	PutTruncate(name string) (CollectionDescriptor, error)

	//PutProperties -> PUT on /_api/collection/{name}/properties
	PutProperties(name string, properties *CollectionPropertyChange) (CollectionDescriptor, error)

	//PutRename -> PUT on /_api/collection/{name}/rename
	PutRename(name string, newName string) (CollectionDescriptor, error)

	//PutRotate -> PUT on /_api/collection/{name}/rotate
	PutRotate(name string) error

	//Delete -> DELETE on /_api/collection/{name}
	Delete(name string) error
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
