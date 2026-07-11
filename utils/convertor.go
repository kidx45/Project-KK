package utils

import (
	"bytes"
	"encoding/json"
	"io"

	db "github.com/kidx45/Project-KK/db/sqlc"
)

func ConvertBuffertoUserType(body *bytes.Buffer) (user db.User, err error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return db.User{}, err
	}
	err = json.Unmarshal(data, &user)
	return
}
