package structs

import "snippetbox/internal/validator"

type SnippetStruct struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Expires int    `json:"expires"`
	validator.Validator
}
