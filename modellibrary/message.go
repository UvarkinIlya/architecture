package modellibrary

import (
	"strings"
	"time"
)

const ImageFeature = "[img]"

type Message struct {
	Text  string
	Time  time.Time
	IsImg bool
}

func NewMessage(text string) Message {
	return Message{
		Text:  text,
		Time:  time.Now(),
		IsImg: isImage(text),
	}
}

func isImage(message string) bool {
	return strings.HasPrefix(message, ImageFeature)
}
