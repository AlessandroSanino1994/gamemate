package gameOwnerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//GameOwnerAction represents a request to enable a game for a user.
type GameOwnerAction struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	GameID       int64  `json:"GameID" xml:"GameID" form:"GameID"`
	UserID       int64  `json:"UserID" xml:"UserID" form:"UserID"`
	Action       bool   `json:"Action" xml:"Action" form:"Action"` //True = Enable, False = Disable (on a GameID)
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *GameOwnerAction) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "GameOwnerAction" || receiver.GameID <= 0 ||
		receiver.UserID <= 0 {
		return errors.New("Invalid Form Submitted " + err.Error())
	}

	return nil
}
