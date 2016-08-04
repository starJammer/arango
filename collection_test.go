package arango

import (
	"net/http"
	"testing"
)

func TestCollectionEndpoint(t *testing.T) {

	var db = getDatabase("_system")
	var ce = db.CollectionEndpoint()

	if ce.Database() == nil {
		t.Fatal("Expected a link back to collections database.")
	}

	if ce.Database().Name() != db.Name() {
		t.Fatalf(
			"Expected db names to match: Actual(%s), Expected(%s)",
			ce.Database().Name(),
			db.Name(),
		)
	}

}

func TestGetCollections(t *testing.T) {
	var ce = getCE("_system")

	collections, err := ce.GetCollections(false)

	if err != nil {
		t.Fatal("Unexpected error during GetCollections: ", err)
	}

	if len(collections) < 1 {
		t.Fatal("Expected at least one collection.")
	}

	for _, coll := range collections {
		if !coll.IsSystem {
			t.Fatal("Expected only system collections in _system db")
		}

		if coll.Id == "" {
			t.Fatal("Expected an id value for the collection.")
		}

		if coll.Name == "" {
			t.Fatal("Expected a name value for the collection.")
		}

		if coll.Status == 0 {
			t.Fatal("Expected a CollectionStatus value for the collection.")
		}

		if coll.Type == 0 {
			t.Fatal("Expected a CollectionType value for the collection.")
		}
	}

	collections, err = ce.GetCollections(true)

	if len(collections) > 0 {
		t.Fatal(
			"Expected no collections when excluding system collections but got: ",
			len(collections),
		)
	}
}

func TestPostBlankNameCollection(t *testing.T) {
	var ce = getCE("_system")

	co, err := ce.PostCollection(&PostCollectionOptions{
		Name: "",
	})

	if co != nil {
		t.Fatal("Unexpected collection descriptor")
	}

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when creating collection with no name or options.",
	)
}

func TestPostBadName(t *testing.T) {
	var ce = getCE("_system")

	co, err := ce.PostCollection(&PostCollectionOptions{
		Name: "_fakesystem",
	})

	if co != nil {
		t.Fatal("Unexpected collection descriptor")
	}

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when creating a collection starting with underscore.",
	)
}

func TestDeleteBlankNameCollection(t *testing.T) {
	var ce = getCE("_system")

	co, err := ce.Delete("")

	if co != nil {
		t.Fatal("Unexpected collection descriptor on delete")
	}

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when deleting blank named collection.",
	)
}

func TestDeleteNonExistingCollection(t *testing.T) {
	var ce = getCE("_system")

	co, err := ce.Delete("non-existent")

	if co != nil {
		t.Fatal("Unexpected collection descriptor on delete")
	}

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when deleting collection that doesn't exist.",
	)
}

func TestPostDeleteCollection(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true

	co, err := ce.PostCollection(opts)

	if co == nil {
		t.Fatal("Expected a collection descriptor after creation: ", err)
	}

	if co.Id == "" {
		t.Fatal("Expected collection descriptor to have the Id")
	}

	if co.Name == "" {
		t.Fatal("Expected collection descriptor to have the Name")
	}

	if !co.WaitForSync {
		t.Fatal("Expected collection description to have a true WaitForSync value")
	}

	if co.Status != LOADED_STATUS {
		t.Fatal("Expected collection description to have a loaded status of", LOADED_STATUS, "but it had", co.Status)
	}

	if err != nil {
		t.Fatal("Unexpected error when creating collection: ", err)
	}

	co2, err := ce.Delete(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when deleting collection: ", err)
	}

	if co.Id != co2.Id {
		t.Fatal("Unexpectedly deleted a different collection: Expected = ", co.Id, "Actual = ", co2.Id)
	}

}

func TestGetBlankNameCollection(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.Get("")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}

}

func TestGetNonExistentCollection(t *testing.T) {
	var ce = getCE("_system")
	d, err := ce.Get("bad-name")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestIncludedInGetCollections(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts)
	defer ce.Delete("test")

	colls, _ := ce.GetCollections(true)

	if found := colls.Find(opts.Name); found == nil || found.Name != opts.Name {
		t.Fatal("Could not find newly created connection.")
	}
}

func TestGetCollection(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts)
	defer ce.Delete("test")

	descriptor, err := ce.Get(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Id == "" {
		t.Fatal("No id present in newly created collection.")
	}

	if descriptor.Name != opts.Name {
		t.Fatalf("Unexpected collection name - Expected(%s) Actual(%s)", opts.Name, descriptor.Name)
	}

	if descriptor.Status != LOADED_STATUS {
		t.Fatalf("Unexpected collection status - Expected(%d) Actual(%d)", LOADED_STATUS, descriptor.Status)
	}

	if descriptor.Type != DOCUMENT_COLLECTION {
		t.Fatalf("Unexpected collection type - Expected(%d) Actual(%d)", DOCUMENT_COLLECTION, descriptor.Type)
	}

	if descriptor.IsSystem != false {
		t.Fatalf("Unexpected IsSystem value - Expected(%t) Actual(%t)", false, descriptor.IsSystem)
	}

}

func TestPostGetEdgeCollection(t *testing.T) {
	var ce = getCE("_system")
	opts := DefaultPostCollectionOptions()
	opts.Type = EDGE_COLLECTION
	opts.Name = "edge-test"

	ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.Get(opts.Name)

	if err != nil {
		t.Fatal("Unexpected result from CollectionEndpoint.Get", err)
	}

	if descriptor.Type != EDGE_COLLECTION {
		t.Fatal("Expected collection to be of type EDGE: ", descriptor.Type)
	}

}

func TestGetPropertiesBlankName(t *testing.T) {

	var ce = getCE("_system")
	d, err := ce.GetProperties("")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching properties of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}

}

func TestGetPropertiesNonExistent(t *testing.T) {

	var ce = getCE("_system")
	d, err := ce.GetProperties("non-existent")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching properties of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}

}

func TestGetCollectionProperties(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true
	opts.DoCompact = false

	_, err := ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetProperties(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.WaitForSync != true {
		t.Fatal("Expected waitForSync to be true.")
	}

	if descriptor.DoCompact != false {
		t.Fatal("Expected doCompact to be false.")
	}

	if descriptor.JournalSize == 0 {
		t.Fatal("Expected positive value for journal size.")
	}

	if descriptor.KeyOptions == nil {
		t.Fatal("Expected some key options.")

		k := descriptor.KeyOptions
		if k.Type == "" {
			t.Fatal("Expected a value for keyoptions type")
		}
	}

	if descriptor.IsVolatile != false {
		t.Fatal("Expected isVolatile to be false.")
	}

	if descriptor.NumberOfShards != 0 {
		t.Fatal("Expected numberOfShards to be 0: ", descriptor.NumberOfShards)
	}

	if len(descriptor.ShardKeys) != 0 {
		t.Fatal("Expected ShardKeys to have a length of zero.: ", descriptor.ShardKeys)
	}

}

func TestGetCollectionCountBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetCount("")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching count of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionCountNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetCount("non-existent")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching count of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionCount(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	_, err := ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetCount(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Count != 0 {
		t.Fatal("Expected to get a count of 0.")
	}

}

func TestGetCollectionFiguresBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetFigures("")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching figures of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionFiguresNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetFigures("non-existent")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching figures of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionFigures(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	_, err := ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetFigures(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Figures == nil {
		t.Fatal("Expected figures to be non-nil.")
	}

}

func TestGetCollectionRevisionBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetRevision("")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching revision of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionRevisionNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetRevision("non-existent")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching revision of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionRevision(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	_, err := ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetRevision(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Revision == "" {
		t.Fatal("Expected revision to be non-blank.")
	}

}

func TestGetCollectionChecksumBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetChecksum(&GetChecksumOptions{
		Name: "",
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching checksum of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionChecksumNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.GetChecksum(&GetChecksumOptions{
		Name: "non-existent",
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when fetching checksum of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestGetCollectionChecksum(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	_, err := ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetChecksum(&GetChecksumOptions{
		Name: opts.Name,
	})

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	//checksum for new collections is 0.
	//Test checksum for collections with docs within
	//later
	if descriptor.Checksum != 0 {
		t.Fatal("Expected checksum to be non-zero.")
	}

}

func TestPutLoadBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutLoad(&PutLoadOptions{
		Name:  "",
		Count: false,
	})

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when putLoading of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func PutLoadNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutLoad(&PutLoadOptions{
		Name:  "non-existent",
		Count: false,
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when putLoading of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutLoad(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, err := ce.PutLoad(&PutLoadOptions{
		Name:  opts.Name,
		Count: false,
	})

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Status != LOADED_STATUS {
		t.Fatal("Expected collection to be in loaded state.")
	}

}

func TestPutUnloadBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutUnload("")

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when putunloading of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func PutUnloadNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutUnload("non-existent")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when putunloading of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutUnload(t *testing.T) {
	var ce = getCE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts)

	defer ce.Delete(opts.Name)

	descriptor, err := ce.PutUnload(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Status == LOADED_STATUS {
		t.Fatal("Expected collection to not be in loaded state: ", descriptor.Status)
	}
}

func TestPutTruncateBlankName(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutTruncate("")

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when putTruncateing of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func PutTruncateNonExistent(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutTruncate("non-existent")

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when putTruncateing of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutTruncate(t *testing.T) {
	var ce = getCE("_system")
	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	d, err := ce.PutTruncate(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when trucating collection : ", err)
	}

	if d == nil {
		t.Fatal("Expected descriptor to be non nil.")
	}

}

func TestPutPropertiesBlankNameNilProps(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutProperties(&PutPropertiesOptions{
		Name: "",
	})

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when putProperties of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutPropertiesNonExistentNilProps(t *testing.T) {
	var ce = getCE("_system")

	d, err := ce.PutProperties(&PutPropertiesOptions{
		Name: "non-existent",
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when putProperties of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutProperties(t *testing.T) {
	var ce = getCE("_system")
	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true
	ce.PostCollection(opts)
	defer ce.Delete(opts.Name)

	descriptor, _ := ce.GetProperties(opts.Name)

	if descriptor.WaitForSync != true {
		t.Fatal("Expected waitforsync to be true upon creation.")
	}

	descriptor, err := ce.PutProperties(&PutPropertiesOptions{
		Name:        opts.Name,
		WaitForSync: false,
	})

	if err != nil {
		t.Fatal("Unexpected error when putting prorties.")
	}

	if descriptor.WaitForSync != false {
		t.Fatal("Expected waitforSync to be false now.")
	}
}

func TestPutRenameBlankNameNilProps(t *testing.T) {
	var ce = getCE("_system")
	var newName = "newtestname"

	d, err := ce.PutRename(&PutRenameOptions{
		OldName: "",
		NewName: newName,
	})

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when PutRename of blank named collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutRenameNonExistentNilProps(t *testing.T) {
	var ce = getCE("_system")
	var newName = "newtestname"

	d, err := ce.PutRename(&PutRenameOptions{
		OldName: "non-existent",
		NewName: newName,
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when PutRename of non-existent collection.",
	)

	if d != nil {
		t.Fatal("Expected collection descriptor to be nil with an error.")
	}
}

func TestPutRename(t *testing.T) {
	var ce = getCE("_system")
	var newName = "newtestname"

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	_, err := ce.PostCollection(opts)

	descriptor, err := ce.PutRename(&PutRenameOptions{
		OldName: opts.Name,
		NewName: newName,
	})

	if err != nil {
		t.Fatal("Error during rename: ", err)
	}

	if descriptor.Name != newName {
		t.Fatal("Collection rename failed: ", descriptor.Name)
	}

	_, err = ce.Get(opts.Name)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected an error when getting properties of old collection.",
	)

	descriptor, err = ce.Get(newName)

	if err != nil {
		t.Fatal("Error getting renamed collection: ", err)
	}

	if descriptor.Name != newName {
		t.Fatal("Getting properties rename failed: ", descriptor.Name)
	}

	co, err := ce.Delete(newName)

	if err != nil {
		t.Fatal("Error deleting newly renamed collection: ", err)
	}

	if co.Id == "" {

	}
}
