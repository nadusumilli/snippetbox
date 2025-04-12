package structs

import (
	"fmt"
	"snippetbox/cmd/web/constants"
	"snippetbox/internal/validator"
)

type SnippetStruct struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // Exclude from form decoding
}

func (s *SnippetStruct) SetValidator(v validator.Validator) {
	s.Validator = v
}

func (s *SnippetStruct) Validate() {
	s.Validator = validator.New(SnippetStruct{})
	s.CheckField(validator.NotBlank(s.Title), "Title", constants.ErrCannotBeBlank)
	s.CheckField(validator.MaxChars(s.Title, 100), "Title", fmt.Sprintf(constants.ErrMaxChars, 100))
	s.CheckField(validator.NotBlank(s.Content), "Content", constants.ErrCannotBeBlank)
	s.CheckField(validator.PermittedValue(s.Expires, 1, 7, 365), "Expires", "This field must equal 1, 7 or 365")
}
