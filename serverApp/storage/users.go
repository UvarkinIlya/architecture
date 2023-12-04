package storage

import (
	"encoding/json"
	"os"

	"architecture/modellibrary"
)

type UserStorage interface {
	GetUser(login string) (user modellibrary.UserWithActions, isExit bool, err error)
}

func (s *StorageImpl) GetUser(login string) (user modellibrary.UserWithActions, isExit bool, err error) {
	_, err = os.OpenFile(s.UserFilepath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return modellibrary.UserWithActions{}, false, err
	}

	messagesBin, err := os.ReadFile(s.UserFilepath)
	if err != nil {
		return modellibrary.UserWithActions{}, false, err
	}

	if len(messagesBin) == 0 {
		return modellibrary.UserWithActions{}, false, nil
	}

	users := make([]modellibrary.UserWithActions, 0)
	err = json.Unmarshal(messagesBin, &users)
	if err != nil {
		return modellibrary.UserWithActions{}, false, err
	}

	user, found := findUser(login, users)
	return user, found, nil
}

func findUser(login string, users []modellibrary.UserWithActions) (user modellibrary.UserWithActions, isFound bool) {
	for _, userDB := range users {
		if userDB.Login != login {
			continue
		}

		return userDB, true
	}

	return modellibrary.UserWithActions{}, false
}
