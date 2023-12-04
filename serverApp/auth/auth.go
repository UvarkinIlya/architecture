package auth

import (
	"encoding/base64"
	"encoding/json"

	"architecture/modellibrary"
	"architecture/serverApp/storage"
)

type Auther interface {
	Auth(user modellibrary.User) (token string, isAuthorised bool, err error)
}

type AutherImpl struct {
	db storage.UserStorage
}

func NewAutherImpl(db storage.Storage) *AutherImpl {
	return &AutherImpl{db: db}
}

func (m *AutherImpl) Auth(user modellibrary.User) (token string, isAuthorised bool, err error) {
	userDB, isExit, err := m.db.GetUser(user.Login)
	if err != nil {
		return "", false, err
	}

	if !isExit || !(userDB.Login == user.Login && userDB.Password == user.Password) {
		return "", false, nil
	}

	return createToken(userDB), true, nil
}

func createToken(user modellibrary.UserWithActions) (token string) {
	userBin, err := json.Marshal(user)
	if err != nil {
		return ""
	}
	token = base64.StdEncoding.EncodeToString(userBin)

	return token
}
