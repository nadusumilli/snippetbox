package structs

import "snippetbox/internal/validator"

type SnippetStruct struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	validator.Validator
}
