package queue

import (
	"encoding/base64"
	"os/user"
	"path"

	"github.com/google/uuid"
)

func BasePath(status MessageStatus) string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	var subDir string
	switch status {
	case MessageStatusPreparing:
		subDir = "preparing"
	case MessageStatusQueued:
		subDir = "queued"
	default:
		panic("Invalid message status")
	}
	return path.Join(user.HomeDir, ".achelous", subDir)
}

func MsgBasePath(status MessageStatus, msgId uuid.UUID) string {
	return path.Join(BasePath(status), msgId.String())
}

func MsgMetaPath(status MessageStatus, msgId uuid.UUID) string {
	return path.Join(MsgBasePath(status, msgId), "meta.json")
}

func MsgTextBodyPath(status MessageStatus, msgId uuid.UUID) string {
	return path.Join(MsgBasePath(status, msgId), "body.txt")
}

func MsgHtmlBodyPath(status MessageStatus, msgId uuid.UUID) string {
	return path.Join(MsgBasePath(status, msgId), "body.html")
}

func AttBasePath(status MessageStatus, msgId uuid.UUID, attId string) string {
	attId = base64.StdEncoding.EncodeToString([]byte(attId))
	return path.Join(MsgBasePath(status, msgId), attId)
}

func AttMetaPath(status MessageStatus, msgId uuid.UUID, attId string) string {
	return path.Join(AttBasePath(status, msgId, attId), "meta")
}

func AttBodyPath(status MessageStatus, msgId uuid.UUID, attId string) string {
	return path.Join(AttBasePath(status, msgId, attId), "body")
}
