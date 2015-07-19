package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/url"
)

type collectionDescriptor struct {
	Idf             string           `json:"id"`
	Namef           string           `json:"name"`
	IsSystemf       bool             `json:"isSystem"`
	Statusf         CollectionStatus `json:"status"`
	Typef           CollectionType   `json:"type"`
	WaitForSyncf    bool             `json:"waitForSync"`
	DoCompactf      bool             `json:"doCompact"`
	JournalSizef    int              `json:"journalSize"`
	KeyOptionsf     *keyOptions      `json:"keyOptions"`
	IsVolatilef     bool             `json:"isVolatile"`
	NumberOfShardsf int              `json:"numberOfShards"`
	ShardKeysf      []string         `json:"shardKeys"`
	Countf          int              `json:"count"`
	Figuresf        *figures         `json:"figures"`
	Revisionf       string           `json:"revision"`
	Checksumf       int              `json:"checksum"`
}

type figures struct {
	Alivef                     *statHolder `json:"alive"`
	Deadf                      *statHolder `json:"dead"`
	Datafilesf                 *statHolder `json:"datafiles"`
	Journalsf                  *statHolder `json:"journals"`
	Compactorsf                *statHolder `json:"compactors"`
	Shapefilesf                *statHolder `json:"shapefiles"`
	Shapesf                    *statHolder `json:"shapes"`
	Attributesf                *statHolder `json:"attributes"`
	Indexesf                   *statHolder `json:"indexes"`
	LastTickf                  string      `json:"lastTick"`
	UncollectedLogfileEntriesf int         `json:"uncollectedLogfileEntries"`
}

func (f *figures) Alive() Alive {
	return f.Alivef
}

func (f *figures) Dead() Dead {
	return f.Deadf
}

func (f *figures) Datafiles() Datafiles {
	return f.Datafilesf
}

func (f *figures) Journals() Journals {
	return f.Journalsf
}

func (f *figures) Compactors() Compactors {
	return f.Compactorsf
}

func (f *figures) Shapefiles() Shapefiles {
	return f.Shapefilesf
}

func (f *figures) Shapes() Shapes {
	return f.Shapesf
}

func (f *figures) Attributes() Attributes {
	return f.Attributesf
}

func (f *figures) Indexes() Indexes {
	return f.Indexesf
}

func (f *figures) MaxTick() string {
	return f.LastTickf
}

func (f *figures) UncollectedLogfileEntries() int {
	return f.UncollectedLogfileEntriesf
}

type statHolder struct {
	Countf    int `json:"count"`
	Sizef     int `json:"size"`
	Deletionf int `json:"deletion"`
	FileSizef int `json:"fileSize"`
}

func (s *statHolder) Count() int {
	return s.Countf
}

func (s *statHolder) Size() int {
	return s.Sizef
}

func (s *statHolder) Deletion() int {
	return s.Deletionf
}

func (s *statHolder) FileSize() int {
	return s.FileSizef
}

type keyOptions struct {
	Typef          string `json:"type,omitempty"`
	AllowUserKeysf bool   `json:"allowUserKeys"`
	Incrementf     int    `json:"increment"`
	Offsetf        int    `json:"offset"`
}

func (k *keyOptions) Type() string {
	return k.Typef
}

func (k *keyOptions) AllowUserKeys() bool {
	return k.AllowUserKeysf
}

func (k *keyOptions) Increment() int {
	return k.Incrementf
}

func (k *keyOptions) Offset() int {
	return k.Offsetf
}

func (c *collectionDescriptor) Id() string {
	return c.Idf
}

func (c *collectionDescriptor) Name() string {
	return c.Namef
}

func (c *collectionDescriptor) IsSystem() bool {
	return c.IsSystemf
}

func (c *collectionDescriptor) Status() CollectionStatus {
	return c.Statusf
}

func (c *collectionDescriptor) Type() CollectionType {
	return c.Typef
}

func (c *collectionDescriptor) WaitForSync() bool {
	return c.WaitForSyncf
}

func (c *collectionDescriptor) DoCompact() bool {
	return c.DoCompactf
}

func (c *collectionDescriptor) JournalSize() int {
	return c.JournalSizef
}

func (c *collectionDescriptor) KeyOptions() KeyOptions {
	return c.KeyOptionsf
}

func (c *collectionDescriptor) IsVolatile() bool {
	return c.IsVolatilef
}

func (c *collectionDescriptor) NumberOfShards() int {
	return c.NumberOfShardsf
}

func (c *collectionDescriptor) ShardKeys() []string {
	return c.ShardKeysf
}

func (c *collectionDescriptor) Count() int {
	return c.Countf
}

func (c *collectionDescriptor) Figures() Figures {
	return c.Figuresf
}

func (c *collectionDescriptor) Revision() string {
	return c.Revisionf
}

func (c *collectionDescriptor) Checksum() int {
	return c.Checksumf
}

type collection struct {
	name     string
	client   gr.Client
	database *database
}

func (c *collection) Name() string {
	return c.name
}

func (c *collection) Database() Database {
	return c.database
}

func (c *collection) Get() (CollectionDescriptor, error) {

	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Get("", nil, descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) GetProperties() (CollectionDescriptor, error) {
	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Get("/properties", nil, descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) GetCount() (CollectionDescriptor, error) {
	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Get("/count", nil, descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) GetFigures() (CollectionDescriptor, error) {
	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Get("/figures", nil, descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) GetRevision() (CollectionDescriptor, error) {
	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Get("/revision", nil, descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) GetChecksum(withRevisions bool, withData bool) (CollectionDescriptor, error) {
	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Get(
		"/checksum",
		url.Values{
			"withRevisions": []string{fmt.Sprintf("%t", withRevisions)},
			"withData":      []string{fmt.Sprintf("%t", withData)},
		},
		descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) PutLoad(includeCount bool) (CollectionDescriptor, error) {

	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Put("/load",
		nil,
		map[string]string{"count": fmt.Sprintf("%t", includeCount)},
		descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) PutUnload() (CollectionDescriptor, error) {

	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Put("/unload",
		nil,
		nil,
		descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) PutTruncate() (CollectionDescriptor, error) {

	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Put("/truncate",
		nil,
		nil,
		descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) PutProperties(properties *CollectionPropertyChange) (CollectionDescriptor, error) {

	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Put("/properties",
		nil,
		properties,
		descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	return descriptor, nil
}

func (c *collection) PutRename(name string) (CollectionDescriptor, error) {

	var descriptor = &collectionDescriptor{}
	var errorResult = &arangoError{}

	h, err := c.client.Put("/rename",
		nil,
		map[string]string{"name": name},
		descriptor, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
		return nil, errorResult
	}

	c.client.BaseUrl().Path = fmt.Sprintf(DatabasePath, c.Database().Name()) + fmt.Sprintf(CollectionPath, name)
	c.name = name
	return descriptor, nil
}

func (c *collection) Delete() error {

	var errorResult = &arangoError{}

	h, err := c.client.Delete("", nil, nil, errorResult)

	if err != nil {
		return err
	}

	if h.StatusCode != 200 {
		return errorResult
	}

	return nil

}
