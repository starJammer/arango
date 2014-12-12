package arango

//Document represents an Arango document.
type FullDocument interface {
	HasArangoId
	HasArangoRev
	HasArangoKey
}

type HasArangoId interface {
	Id() string
	SetId(string)
}

type HasArangoRev interface {
	Rev() string
	SetRev(string)
}

type HasArangoKey interface {
	Key() string
	SetKey(string)
}

//DocumentImplementation is an embeddable type that
//you can use that already implements the Document interface.
type DocumentImplementation struct {
	ArangoId  string `json:"_id,omitempty"`
	ArangoRev string `json:"_rev,omitempty"`
	ArangoKey string `json:"_key,omitempty"`
}

func (d *DocumentImplementation) Id() string {
	return d.ArangoId
}

func (d *DocumentImplementation) SetId(id string) {
	d.ArangoId = id
}

func (d *DocumentImplementation) Rev() string {
	return d.ArangoRev
}

func (d *DocumentImplementation) SetRev(rev string) {
	d.ArangoRev = rev
}

func (d *DocumentImplementation) Key() string {
	return d.ArangoKey
}

func (d *DocumentImplementation) SetKey(key string) {
	d.ArangoKey = key
}

//GetOptions are used when fetching documents
//Read the GET /_api/document/{document-handle} info
type GetOptions struct {
	IfNoneMatch string
	IfMatch     string
}

//Save options represent options available when calling the Post /_api/document/ endpoint
type SaveOptions struct {
	//The collection to save the item to. Irrelevant when called from a Collection struct
	Collection string

	//CreateCollection specifies if the collection should be created at the same time as this document is being saved.
	//Irrelevant if called from
	CreateCollection bool

	//Wait until document has been synced to disk.
	WaitForSync bool
}

type ReplaceOptions struct {
	WaitForSync bool
	Rev         string
	Policy      string
	IfMatch     string
}

//DefaultReplaceOptions returns options with default values according to arango
func DefaultReplaceOptions() *ReplaceOptions {
	return &ReplaceOptions{
		WaitForSync: false,
		Rev:         "",
		Policy:      "error",
		IfMatch:     "",
	}
}

type UpdateOptions struct {
	KeepNull    bool
	MergeArrays bool

	WaitForSync bool
	Rev         string
	Policy      string
	IfMatch     string
}

//DefaultUpdateOptions returns options with default values according to arango
func DefaultUpdateOptions() *UpdateOptions {

	return &UpdateOptions{
		KeepNull:    true,
		MergeArrays: true,
		WaitForSync: false,
		Rev:         "",
		Policy:      "error",
		IfMatch:     "",
	}
}
