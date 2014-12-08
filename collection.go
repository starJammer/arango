package arango

import (
	"fmt"
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

func (c *Collection) Id() string {
	return c.json.Id
}

func (c *Collection) Name() string {
	return c.json.Name
}

func (c *Collection) Status() int {
	return c.json.Status
}

func (c *Collection) Type() int {
	return c.json.Type
}

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
