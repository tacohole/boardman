package schema

type Source struct {
	ID           int
	Name         string
	Url          string
	Endpoint     string
	AuthRequired bool
	AuthType     string
	ResponseId   int
	Response     Response
	SchemaId     int
	Schema       Schema
}

type Response struct {
	ID              int
	Name            string
	SourceId        int
	SchemaId        int
	ContentType     string
	DefaultPageSize int
	PageIndexVar    string
	NextPageVar     string
}

type Schema struct {
	ID     int
	Name   string
	Fields []string
}
