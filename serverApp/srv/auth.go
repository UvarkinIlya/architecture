package srv

import (
	"encoding/json"
	"fmt"
	"net/http"

	"architecture/logger"
	"architecture/modellibrary"
)

func (s *ServerImpl) auth(writer http.ResponseWriter, request *http.Request) {
	var user modellibrary.User

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		logger.Error("Failed decode user due to error:%s", err)
		_, _ = fmt.Fprintf(writer, err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, isAuthorised, err := s.auther.Auth(user)
	if err != nil {
		_, _ = fmt.Fprintf(writer, err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isAuthorised {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = fmt.Fprintf(writer, token)
	if err != nil {
		_, _ = fmt.Fprintf(writer, err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
