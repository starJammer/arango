package arango

type CollectionType int

//Collection types
const (
	DOCUMENT_COLLECTION CollectionType = 2
	EDGE_COLLECTION     CollectionType = 3
)

type CollectionStatus int

//Collection statuses are all the states a
//collection can be in.
const (
	NEW_BORN_STATUS       CollectionStatus = 1
	UNLOADED_STATUS       CollectionStatus = 2
	LOADED_STATUS         CollectionStatus = 3
	BEING_UNLOADED_STATUS CollectionStatus = 4
	DELETED_STATUS        CollectionStatus = 5
	LOADING_STATUS        CollectionStatus = 6
)

//Policy represents the different
//policy types that arango supports when
//putting or patching a document
type Policy string

//The two policy types are error and last
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
