package schema

type Source struct {
	ID              int      `csv:"ID"json:"_id"`
	Name            string   `csv:"Name"json:"name"`
	Url             string   `csv:"URL"json:"url"`
	AuthRequired    bool     `csv:"Auth Required"json:"authRequired"`
	AuthType        string   `csv:"Auth Type"json:"authType"`
	ContentType     string   `csv:"Content Type"json:"contentType"`
	DefaultPageSize int      `csv:"Page Size"json:"pageSize"`
	PageIndexVar    string   `csv:"Page Index"json:"pageIndex"`
	NextPageVar     string   `csv:"Next Page"json:"nextPage"`
	DbSchemaId      int      `csv:"DB Schema ID"json:"dbSchemaId"`
	DbSchema        DbSchema `csv:"DB Schema Name"json:"dbSchema"`
}

type DbSchema struct {
	ID     int      `csv:"ID"json:"_id"`
	Name   string   `csv:"Name"json:"name"`
	Fields []string `csv:"Fields"json:"fields"`
}
