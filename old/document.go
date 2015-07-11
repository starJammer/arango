package old

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

type ArangoEdge interface {
	From() string
	SetFrom(string)

	To() string
	SetTo(string)
}

//DocumentImplementation is an embeddable type that
//you can use that already implements the Document interfaces.
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

//EdgeImplementation is an embeddable type that
//you can use that already implements the edge interfaces.
type EdgeImplementation struct {
	DocumentImplementation
	ArangoFrom string `json:"_from,omitempty"`
	ArangoTo   string `json:"_to,omitempty"`
}

func (e *EdgeImplementation) From() string {
	return e.ArangoFrom
}

func (e *EdgeImplementation) SetFrom(from string) {
	e.ArangoFrom = from
}

func (e *EdgeImplementation) To() string {
	return e.ArangoTo
}

func (e *EdgeImplementation) SetTo(to string) {
	e.ArangoTo = to
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
	//Irrelevant if called from a collection struct
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
//You don't really need this unless you plan on using the Rev/IfMatch or
//WaitForSync options. In that case, it helps to have this method to have
//the defaults that arango will use. That way you only change what you need.
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
//You don't really need this unless you plan on using the Rev/IfMatch or
//other options. In that case, it helps to have this method to have
//the defaults that arango will use. That way you only change what you need.
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

type DeleteOptions ReplaceOptions

func DefaultDeleteOptions() *DeleteOptions {
	return &DeleteOptions{
		WaitForSync: false,
		Rev:         "",
		Policy:      "error",
		IfMatch:     "",
	}
}
