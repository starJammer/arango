package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
)

type CollectionDescriptors []CollectionDescriptor

func (c CollectionDescriptors) Find(name string) *CollectionDescriptor {
	for _, d := range c {
		if d.Name == name {
			return &d
		}
	}
	return nil
}

type CollectionDescriptor struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	IsSystem       bool             `json:"isSystem"`
	Status         CollectionStatus `json:"status"`
	Type           CollectionType   `json:"type"`
	WaitForSync    bool             `json:"waitForSync"`
	DoCompact      bool             `json:"doCompact"`
	JournalSize    int              `json:"journalSize"`
	KeyOptions     *KeyOptions      `json:"keyOptions"`
	IsVolatile     bool             `json:"isVolatile"`
	NumberOfShards int              `json:"numberOfShards"`
	ShardKeys      []string         `json:"shardKeys"`
	Count          int              `json:"count"`
	Figures        *Figures         `json:"figures"`
	Revision       string           `json:"revision"`
	Checksum       int              `json:"checksum"`
}

type Figures struct {
	Alive                     Alive      `json:"alive"`
	Dead                      Dead       `json:"dead"`
	Datafiles                 Datafiles  `json:"datafiles"`
	Journals                  Journals   `json:"journals"`
	Compactors                Compactors `json:"compactors"`
	Shapefiles                Shapefiles `json:"shapefiles"`
	Shapes                    Shapes     `json:"shapes"`
	Attributes                Attributes `json:"attributes"`
	Indexes                   Indexes    `json:"indexes"`
	LastTick                  string     `json:"lastTick"`
	UncollectedLogfileEntries int        `json:"uncollectedLogfileEntries"`
}

type StatHolder struct {
	Count    int `json:"count"`
	Size     int `json:"size"`
	Deletion int `json:"deletion"`
	FileSize int `json:"fileSize"`
}

type Alive StatHolder
type Dead StatHolder
type Datafiles StatHolder
type Journals StatHolder
type Compactors StatHolder
type Shapefiles StatHolder
type Shapes StatHolder
type Attributes StatHolder
type Indexes StatHolder

type KeyOptions struct {
	Type          string `json:"type,omitempty"`
	AllowUserKeys bool   `json:"allowUserKeys"`
	Increment     int    `json:"increment"`
	Offset        int    `json:"offset"`
}

type CollectionEndpoint struct {
	client   *gr.Client
	database *Database
}

//Database gets the related database endpoint
//for this collection endpoint
func (c *CollectionEndpoint) Database() *Database {
	return c.database
}

//GetCollections -> GET on  /_api/collection
func (c *CollectionEndpoint) GetCollections(excludeSystemCollections bool) (CollectionDescriptors, error) {

	var result struct {
		Collections []CollectionDescriptor `json:"collections"`
	}

	var errorResult = ArangoError{}

	h, err := c.client.Get(
		"",
		nil,
		url.Values{"excludeSystem": []string{fmt.Sprintf("%t", excludeSystemCollections)}},
		gr.UnmarshalMap{
			http.StatusOK: &result,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return result.Collections, nil
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

//PostCollection -> POST on /_api/collection
func (c *CollectionEndpoint) PostCollection(name string, options *PostCollectionOptions) error {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	if options == nil {
		options = DefaultPostCollectionOptions()
	}
	options.Name = name

	h, err := c.client.Post(
		"",
		nil,
		nil,
		options,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil
}

//Get -> GET on /_api/collection/{name}
func (c *CollectionEndpoint) Get(name string) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	//if name is blank, this is like calling GetCollections(false)
	//Instead, we make name = "-" so we get an appropriate error
	//when name is blank.
	if name == "" {
		name = "-"
	}

	h, err := c.client.Get(
		fmt.Sprintf("/%s", name),
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:       descriptor,
			http.StatusNotFound: &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//GetProperties -> GET on /_api/collection/{name}/properties
func (c *CollectionEndpoint) GetProperties(name string) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Get(
		fmt.Sprintf("/%s/properties", name),
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:       descriptor,
			http.StatusNotFound: &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//GetCount -> GET on /_api/collection/{name}/count
func (c *CollectionEndpoint) GetCount(name string) (*CollectionDescriptor, error) {
	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Get(
		fmt.Sprintf("/%s/count", name),
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//GetFigures -> GET on /_api/collection/{name}/figures
func (c *CollectionEndpoint) GetFigures(name string) (*CollectionDescriptor, error) {
	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Get(
		fmt.Sprintf("/%s/figures", name),
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//GetRevision -> GET on /_api/collection/{name}/revision
func (c *CollectionEndpoint) GetRevision(name string) (*CollectionDescriptor, error) {
	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Get(
		fmt.Sprintf("/%s/revision", name),
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

type GetChecksumOptions struct {
	WithRevisions bool
	WithData      bool
}

//GetChecksum -> GET on /_api/collection/{name}/checksum
func (c *CollectionEndpoint) GetChecksum(name string, opts *GetChecksumOptions) (*CollectionDescriptor, error) {
	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	var query url.Values
	if opts != nil {
		query = make(url.Values)
		query.Add("withRevisions", fmt.Sprintf("%t", opts.WithRevisions))
		query.Add("withData", fmt.Sprintf("%t", opts.WithData))
	}

	h, err := c.client.Get(
		fmt.Sprintf("/%s/checksum", name),
		nil,
		query,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//PutLoad -> PUT on /_api/collection/{name}/load
func (c *CollectionEndpoint) PutLoad(name string, includeCount bool) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Put(
		fmt.Sprintf("/%s/load", name),
		nil,
		nil,
		map[string]string{"count": fmt.Sprintf("%t", includeCount)},
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//PutUnload -> PUT on /_api/collection/{name}/unload
func (c *CollectionEndpoint) PutUnload(name string) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Put(
		fmt.Sprintf("/%s/unload", name),
		nil,
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//PutTruncate -> PUT on /_api/collection/{name}/truncate
func (c *CollectionEndpoint) PutTruncate(name string) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Put(
		fmt.Sprintf("/%s/truncate", name),
		nil,
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

type PutPropertiesOptions struct {
	//If the WaitForSync of your collection is true,
	//make sure to set this to true if you're only
	//setting JournalSize, otherwise you will
	//set it to false by mistake.
	WaitForSync bool `json:"waitForSync"`
	//Omitempty because arango expects at least 1048576 bytes (1MB)
	//If you set it to >0 it will be sent in the request but you
	//might get an error if it doesn't meet the minimum requirement.
	JournalSize int `json:"journalSize,omitempty"`
}

//PutProperties -> PUT on /_api/collection/{name}/properties
func (c *CollectionEndpoint) PutProperties(name string, properties PutPropertiesOptions) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Put(
		fmt.Sprintf("/%s/properties", name),
		nil,
		nil,
		&properties,
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//PutRename -> PUT on /_api/collection/{name}/rename
func (c *CollectionEndpoint) PutRename(name string, newName string) (*CollectionDescriptor, error) {

	var descriptor = &CollectionDescriptor{}
	var errorResult = ArangoError{}

	h, err := c.client.Put(
		fmt.Sprintf("/%s/rename", name),
		nil,
		nil,
		map[string]string{"name": newName},
		gr.UnmarshalMap{
			http.StatusOK:         descriptor,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return descriptor, nil
}

//PutRotate -> PUT on /_api/collection/{name}/rotate
func (c *CollectionEndpoint) PutRotate(name string) error {

	var errorResult = ArangoError{}

	h, err := c.client.Put(
		fmt.Sprintf("/%s/rotate", name),
		nil,
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil
}

//Delete -> DELETE on /_api/collection/{name}
func (c *CollectionEndpoint) Delete(name string) error {

	var errorResult = ArangoError{}

	h, err := c.client.Delete(
		fmt.Sprintf("/%s", name),
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil

}
