package arango

import (
	gr "github.com/starJammer/grestclient"
)

type SimpleEndpoint struct {
	client   *gr.Client
	database *Database
}

//Database gets the related database endpoint
//for this collection endpoint
func (s *SimpleEndpoint) Database() *Database {
	return nil

}

type ByExampleOptions struct {
	Skip      int `json:"skip,omitempty"`
	Limit     int `json:"limit,omitempty"`
	BatchSize int `json:"batchSize,omitempty"`
}

func (s *SimpleEndpoint) ByExample(collection string, example interface{}, opts *ByExampleOptions) *Cursor {

	return nil
}
