package web

import (
	//"github.com/rivo/users"
	"log"
	"time"
)

type ExampleUser struct {
	Id             string
	email          string
	passwordHash   []byte
	state          int
	verificationID string
	vidCreated     time.Time
	passwordToken  string
	tokenCreated   time.Time
}

func (u *ExampleUser) GetID() interface{} {
	return u.Id
}

func (u *ExampleUser) SetID(id interface{}) {
	u.Id = id.(string)
}

func (u *ExampleUser) SetState(state int) {
	u.state = state
}

func (u *ExampleUser) GetState() int {
	log.Print("GetState")
	return 1
	//return u.state
}

func (u *ExampleUser) SetEmail(email string) {
	u.email = email
}

func (u *ExampleUser) GetEmail() string {
	log.Print("GetEmail")
	return u.email
}

func (u *ExampleUser) SetPasswordHash(hash []byte) {
	u.passwordHash = hash
}

func (u *ExampleUser) GetPasswordHash() []byte {
	log.Print("GetPasswordHash")
	return u.passwordHash
}

func (u *ExampleUser) SetVerificationID(id string, created time.Time) {
	u.verificationID = id
	u.vidCreated = created
}

func (u *ExampleUser) GetVerificationID() (string, time.Time) {
	return u.verificationID, u.vidCreated
}

func (u *ExampleUser) SetPasswordToken(id string, created time.Time) {
	u.passwordToken = id
	u.tokenCreated = created
}

func (u *ExampleUser) GetPasswordToken() (string, time.Time) {
	return u.passwordToken, u.tokenCreated
}

func (u *ExampleUser) GetRoles() []string {
	return nil
}

//func (u *ExampleUser) LoadUserByEmail(email string) (*ExampleUser, error) {
//
//}
//func NewUser() users.User{
//
//	return users.User{}
//}
