package queue

import (
	"encoding/base64"
	"os/user"
	"path"

	"github.com/google/uuid"
)

func BasePath() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(user.HomeDir, ".achelous/queue/")
}

func MsgBasePath(msgId uuid.UUID) string {
	return path.Join(BasePath(), msgId.String())
}

func MsgMetaPath(msgId uuid.UUID) string {
	return path.Join(MsgBasePath(msgId), "meta.json")
}

func MsgTextBodyPath(msgId uuid.UUID) string {
	return path.Join(MsgBasePath(msgId), "body.txt")
}

func MsgHtmlBodyPath(msgId uuid.UUID) string {
	return path.Join(MsgBasePath(msgId), "body.html")
}

func AttBasePath(msgId uuid.UUID, attId string) string {
	attId = base64.StdEncoding.EncodeToString([]byte(attId))
	return path.Join(MsgBasePath(msgId), attId)
}

func AttMetaPath(msgId uuid.UUID, attId string) string {
	return path.Join(AttBasePath(msgId, attId), "meta")
}

func AttBodyPath(msgId uuid.UUID, attId string) string {
	return path.Join(AttBasePath(msgId, attId), "body")
}
