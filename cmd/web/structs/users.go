package structs

import (
	"fmt"
	"snippetbox/cmd/web/constants"
	"snippetbox/internal/validator"
)

type UserStruct struct {
	Name                string     `form:"name"` // Exclude from form decoding
	Username            string     `form:"username"`
	Password            string     `form:"password"`
	Email               string     `form:"email"`
	validator.Validator `form:"-"` // Exclude from form decoding
}

// Create a new userLoginForm struct.
type UserLogin struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (u *UserStruct) Validate() {
	u.Validator = validator.New(UserStruct{})
	u.Validator.CheckField(validator.NotBlank(u.Name), "Name", constants.ErrCannotBeBlank)
	u.Validator.CheckField(validator.NotBlank(u.Username), "Username", constants.ErrCannotBeBlank)
	u.Validator.CheckField(validator.NotBlank(u.Password), "Password", constants.ErrCannotBeBlank)
	u.Validator.CheckField(validator.NotBlank(u.Email), "Email", constants.ErrCannotBeBlank)
	u.Validator.CheckField(validator.Matches(u.Email, validator.EmailRX), "Email", constants.ErrInvalidEmail)
	u.Validator.CheckField(validator.MinChars(u.Password, 8), "Password", fmt.Sprintf(constants.ErrMinChars, 8))
}

func (u *UserLogin) Validate() {
	u.Validator = validator.New(UserLogin{})
	u.Validator.CheckField(validator.NotBlank(u.Email), "Email", constants.ErrCannotBeBlank)
	u.Validator.CheckField(validator.Matches(u.Email, validator.EmailRX), "Email", constants.ErrInvalidEmail)
	u.Validator.CheckField(validator.NotBlank(u.Password), "Password", constants.ErrCannotBeBlank)
}
