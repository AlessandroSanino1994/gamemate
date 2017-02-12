package developerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//AddToken represents a request to add a token for an app of the developer.
type AddToken struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
func (receiver *AddToken) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")
	if receiver.Type != "AddToken" || receiver.API_Token == "" || receiver.SessionToken == "" {
		err = errors.New("Invalid Form Submitted")
	}

	return err
}