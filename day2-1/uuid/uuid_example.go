package uuid

import (
    "strings"

    "github.com/pborman/uuid"
)

func GenerateUUID() string {
    uuidWithHyphen := uuid.NewRandom()
    uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
    return uuid
}
