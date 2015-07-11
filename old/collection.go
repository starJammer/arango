package old

import (
	"fmt"
	"strings"
)

//Collection types
//
//See arango manual or rest api docs for what these might mean
const (
	DOCUMENT_COLLECTION = 2
	EDGE_COLLECTION     = 3
)

//Collection statuses
//
//See arango manual or rest api docs for what these might mean
const (
	NEW_BORN_STATUS       = 1
	UNLOADED_STATUS       = 2
	LOADED_STATUS         = 3
	BEING_UNLOADED_STATUS = 4
	DELETED_STATUS        = 5
)

//Collection represents a collection from arangodb
//Don't instantiate this yourself. Use db.Collection
//to get the one you want.
type Collection struct {
	db   *Database
	json *collectionResult
}

//CollectionCreationOptions represent options when creating a new collection.
//Look at the documentation for the POST to /_api/collection
//for information on the defaults, optionals, and required
//attributes.
type CollectionCreationOptions struct {
	WaitForSync    bool        `json:"waitForSync,omitempty"`
	DoCompact      bool        `json:"doCompact"`
	JournalSize    int         `json:"journalSize,omitempty"`
	IsSystem       bool        `json:"isSystem,omitempty"`
	IsVolatile     bool        `json:"isVolatile,omitempty"`
	KeyOptions     *KeyOptions `json:"keyOptions,omitempty"`
	Type           int         `json:"type,omitempty"`
	NumberOfShards int         `json:"numberOfShards,omitempty"`
	ShardKeys      []string    `json:"shardKeys,omitempty"`
	//You can set name manually but it will be overriden with
	//whatever the create methods are called with so don't bother.
	Name string `json:"name"`
}

//KeyOptions stores information about how a collection's key is configured.
//It is used during collection creation to specify how the new collection's
//key should be setup.
//
//It is also used for existing collections so you know how the collection's
//key is configured.
//
//If you've fetched KeyOptions by calling c.Properties then
//treat these as read only values.
//You changing them yourself won't do anything special.
type KeyOptions struct {
	Type          string `json:"type,omitempty"`
	AllowUserKeys bool   `json:"allowUserKeys"`
	Increment     int    `json:"increment"`
	Offset        int    `json:"offset"`
}

type collectionResult struct {
	//Fetched during normal call to db.Collection
	Id       string
	Name     string
	Status   int
	Type     int
	IsSystem bool

	//Only populated when you call c.Properties
	//Otherwise these'll be blank
	WaitForSync    bool
	DoCompact      bool
	JournalSize    int
	IsVolatile     bool
	NumberOfShards int
	ShardKeys      []string
	KeyOptions     *KeyOptions

	ArangoError
}

//DefaultCollectionOptions creates a default set of collection options
//If you will always be using the defaults then just use the Create
//method as it uses the defaults.
func DefaultCollectionOptions() CollectionCreationOptions {
	return CollectionCreationOptions{
		DoCompact: true,
	}
}

//Id returns the id of the collection.
func (c *Collection) Id() string {
	return c.json.Id
}

//Name returns the name of the collection.
func (c *Collection) Name() string {
	return c.json.Name
}

//Status returns the status of the collection.
//The result is cached. Call Properties() to refresh.
//The status can be any of the constants defined above.
//const (
//NEW_BORN_STATUS       = 1
//UNLOADED_STATUS       = 2
//LOADED_STATUS         = 3
//BEING_UNLOADED_STATUS = 4
//DELETED_STATUS        = 5
//)
func (c *Collection) Status() int {
	return c.json.Status
}

//Type returns the type of the collection.
//Either it is a document or an edge collection
//The type is returned as one of the two constants:
//  DOCUMENT_COLLECTION = 2
//  EDGE_COLLECTION     = 3
func (c *Collection) Type() int {
	return c.json.Type
}

//IsSystem returns whether the collection is a system collection or not.
//System collections typically start with an underscore like _system
func (c *Collection) IsSystem() bool {
	return c.json.IsSystem
}

//WaitForSync will only have an accurate answer if you call c.Properties first
func (c *Collection) WaitForSync() bool {
	return c.json.WaitForSync
}

//DoCompact will only have an accurate answer if you call c.Properties first
func (c *Collection) DoCompact() bool {
	return c.json.DoCompact
}

//JournalSize will only have an accurate answer if you call c.Properties first
func (c *Collection) JournalSize() int {
	return c.json.JournalSize
}

//IsVolatile will only have an accurate answer if you call c.Properties first
func (c *Collection) IsVolatile() bool {
	return c.json.IsVolatile
}

//NumberOfShards will only have an accurate answer if you call c.Properties first
func (c *Collection) NumberOfShards() int {
	return c.json.NumberOfShards
}

//ShardKeys will only have an accurate answer if you call c.Properties first
func (c *Collection) ShardKeys() []string {
	return c.json.ShardKeys
}

//KeyOptions will only have an accurate answer if you call c.Properties first
func (c *Collection) KeyOptions() *KeyOptions {
	return c.json.KeyOptions
}

//Properties fetches additional properties of the collection.
//It queries the /_api/collection/{collection-name}/properties
//endpoint and causes the collection to switch to the loaded
//state if it was unloaded before.
//Use this especially if you want to update some of the properties
//like the Status().
func (c *Collection) Properties() error {

	db := c.db

	var e ArangoError

	endpoint := fmt.Sprintf("%s/collection/%s/properties", db.serverUrl.String(), c.Name())

	response, err := db.session.Get(endpoint, nil, c.json, &e)

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

//Drop deletes the collection from the database
//DO NOT expect the collection to work after dropping it.
//Calling any further methods on it will result in
//unexpected behavior
func (c *Collection) Drop() error {
	return c.db.DropCollection(c.Name())
}

//Save creates a document in the collection.
//Uses the POST /_api/document endpoint
//If your document embeds the DocumentImplementation type
//or it has fields to hold the Id, Rev, and Key fields
//from arango, then it will be populated with the Id, Rev, Key
//fields during the json.Unmarshal call
func (c *Collection) Save(document interface{}) error {
	return c.db.SaveDocumentWithOptions(document, &SaveOptions{
		Collection:       c.Name(),
		CreateCollection: false,
		WaitForSync:      false,
	})
}

//SaveWithOptions lets you save a document but lets you specify some options
//See the POST /_api/document endpoint for more info.
func (c *Collection) SaveWithOptions(document interface{}, options *SaveOptions) error {

	options.Collection = c.Name()
	options.CreateCollection = false

	return c.db.SaveDocumentWithOptions(document, options)
}

//SaveEdge creates a new edge using pointing from "from" to "to".
//Will probably result in an error if arango determines that this collection is not an edge collection.
func (c *Collection) SaveEdge(from, to, edge interface{}) error {
	return c.db.SaveEdgeWithOptions(from, to, edge, &SaveOptions{
		Collection:       c.Name(),
		CreateCollection: false,
		WaitForSync:      false,
	})
}

//SaveEdgeWithOptions creates a new edge using pointing from "from" to "to" and allows you to specify more options.
//Will probably result in an error if arango determines that this collection is not an edge collection.
func (c *Collection) SaveEdgeWithOptions(from, to, edge interface{}, options *SaveOptions) error {
	options.Collection = c.Name()
	options.CreateCollection = false

	return c.db.SaveEdgeWithOptions(from, to, edge, options)
}

//Document will fetch the document associated with the documentHandle.
//An error will be returned if anything goes wrong. Otherwise, document
//be populated with the document values from arango. json.Unmarhshal
//is used under the hood so make sure you understand how string literals
//tags work with the encoding/json package.
//This uses the GET /_api/document/{document-handle} endpoint.
//The method provides a extra bit of functionality in that it won't let
//you fetch an document NOT from this collection. For example,
//you can't request 'users/123413' when the collection name is 'places'.
//If you want to look up an arbitrary document then use db.Document()
func (c *Collection) Document(documentHandle interface{},
	document interface{}) error {
	return c.DocumentWithOptions(documentHandle, document, nil)
}

func (c *Collection) DocumentWithOptions(documentHandle interface{},
	document interface{},
	options *GetOptions) error {

	documentHandle, ok := c.crossCollectionCheck(documentHandle)
	if ok {
		return c.db.DocumentWithOptions(documentHandle, document, options)
	} else {
		return newError(fmt.Sprintf("Cross collection requests are not permitted.", documentHandle, c.Name()))
	}
}

func (c *Collection) Edge(documentHandle interface{},
	edge interface{}) error {
	return c.EdgeWithOptions(documentHandle, edge, nil)
}

func (c *Collection) EdgeWithOptions(documentHandle interface{},
	edge interface{},
	options *GetOptions) error {

	documentHandle, ok := c.crossCollectionCheck(documentHandle)
	if ok {
		return c.db.EdgeWithOptions(documentHandle, edge, options)
	} else {
		return newError(fmt.Sprintf("Cross collection requests are not permitted.", documentHandle, c.Name()))
	}
}

func (c *Collection) Replace(documentHandle interface{},
	document interface{}) error {
	return c.ReplaceWithOptions(documentHandle, document, nil)
}

func (c *Collection) ReplaceWithOptions(documentHandle interface{},
	document interface{},
	options *ReplaceOptions) error {
	documentHandle, ok := c.crossCollectionCheck(documentHandle)
	if ok {
		return c.db.ReplaceDocumentWithOptions(documentHandle, document, options)
	} else {
		return newError(fmt.Sprintf("Cross collection requests are not permitted.", documentHandle, c.Name()))
	}
}

func (c *Collection) ReplaceEdge(documentHandle interface{},
	edge interface{}) error {
	return c.ReplaceEdgeWithOptions(documentHandle, edge, nil)
}

func (c *Collection) ReplaceEdgeWithOptions(documentHandle interface{},
	edge interface{},
	options *ReplaceOptions) error {
	documentHandle, ok := c.crossCollectionCheck(documentHandle)
	if ok {
		return c.db.ReplaceEdgeWithOptions(documentHandle, edge, options)
	} else {
		return newError(fmt.Sprintf("Cross collection requests are not permitted.", documentHandle, c.Name()))
	}
}

func (c *Collection) Update(documentHandle interface{},
	document interface{}) error {
	return c.UpdateWithOptions(documentHandle, document, nil)
}

func (c *Collection) UpdateWithOptions(documentHandle interface{},
	document interface{},
	options *UpdateOptions) error {
	documentHandle, ok := c.crossCollectionCheck(documentHandle)
	if ok {
		return c.db.UpdateDocumentWithOptions(documentHandle, document, options)
	} else {
		return newError(fmt.Sprintf("Cross collection requests are not permitted.", documentHandle, c.Name()))
	}
}

func (c *Collection) UpdateEdge(documentHandle interface{},
	edge interface{}) error {
	return c.UpdateEdgeWithOptions(documentHandle, edge, nil)
}

func (c *Collection) UpdateEdgeWithOptions(documentHandle interface{},
	edge interface{},
	options *UpdateOptions) error {
	documentHandle, ok := c.crossCollectionCheck(documentHandle)
	if ok {
		return c.db.UpdateDocumentWithOptions(documentHandle, edge, options)
	} else {
		return newError(fmt.Sprintf("Cross collection requests are not permitted.", documentHandle, c.Name()))
	}
}

func (c *Collection) ByExample(example interface{}) (*Cursor, error) {
	return c.db.ByExampleQuery(&ByExampleQuery{
		Collection: c.Name(),
		Example:    example,
	})
}

func (c *Collection) ByExampleQuery(query *ByExampleQuery) (*Cursor, error) {
	if query == nil {
		query = &ByExampleQuery{
			Example: &struct{}{},
		}
	}
	query.Collection = c.Name()
	return c.db.ByExampleQuery(query)
}

func (c *Collection) FirstExample(example, document interface{}) error {
	return c.db.FirstExample(&FirstExampleQuery{
		Collection: c.Name(),
		Example:    example,
	}, document)
}

func (c *Collection) crossCollectionCheck(documentHandle interface{}) (interface{}, bool) {

	switch id := documentHandle.(type) {
	case string:
		idParts := strings.Split(id, "/")
		if len(idParts) == 1 {
			return c.Name() + "/" + idParts[0], true
		} else if len(idParts) == 2 {
			if idParts[0] == c.Name() {
				return id, true
			}
		}
	case HasArangoId:
		idParts := strings.Split(id.Id(), "/")
		if len(idParts) == 2 {
			if idParts[0] == c.Name() {
				return documentHandle, true
			}
		}
	case HasArangoKey:
		return c.Name() + "/" + id.Key(), true
	}

	return "", false
}
