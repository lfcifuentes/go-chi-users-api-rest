package structures

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID
	Username  string `json:"username"`
	Name      string `json:"name"`
	Last_name string `json:"last_name"`
}

type Response struct {
	Status  int    `json:"status"`
	Data    User   `json:"data"`
	Message string `json:"message"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id, err = uuid.NewV1()

	if err != nil {
		return errors.New("can't save invalid data")
	}
	return
}

func (u *User) IsValid() bool {
	if u.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return false
	}
	return true
}
