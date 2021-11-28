package schema

type Data interface {
	StructKeysToString() string
}

type Source struct {
	ID           int
	Name         string
	Url          string
	Endpoint     string
	AuthRequired bool
	AuthType     string
	ResponseId   int
	Response     Response
	DbSchemaId   int
	DbSchema     DbSchema
}

type Response struct {
	ID              int
	Name            string
	SourceId        int
	DbSchemaId      int
	ContentType     string
	DefaultPageSize int
	PageIndexVar    string
	NextPageVar     string
}

type DbSchema struct {
	ID     int
	Name   string
	Fields []string
}
