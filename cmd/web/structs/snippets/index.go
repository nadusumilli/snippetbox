package structs

type SnippetStruct struct {
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Expires     int               `json:"expires"`
	FieldErrors map[string]string `json:"field_errors"`
}
