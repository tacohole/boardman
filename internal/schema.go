package schema

type Source struct {
	ID              int
	Name            string
	Url             string
	AuthRequired    bool
	AuthType        string
	ContentType     string
	DefaultPageSize int
	PageIndexVar    string
	NextPageVar     string
	DbSchemaId      int
	DbSchema        DbSchema
}

type DbSchema struct {
	ID     int
	Name   string
	Fields []string
}
