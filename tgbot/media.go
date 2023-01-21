package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Media is a generic type for all kinds of media that includes File.
type Media interface {
	// MediaType returns string-represented media type.
	MediaType() string

	// MediaFile returns a pointer to the media file.
	MediaFile() *File
}

// InputMedia represents a composite InputMedia struct that is
// used by Telebot in sending and editing media methods.
type InputMedia struct {
	Type    string `json:"type"`
	Caption string `json:"caption"`
}

// Document object represents a general file (as opposed to Photo or Audio).
// Telegram users can send files of any type of up to 1.5 GB in size.
type Document struct {
	File
	tgbotapi.Document
	Caption string
}

func NewDocument(filename string, data []byte) *Document {
	return &Document{
		File: File{
			filename: filename,
			data:     data,
		},
	}
}

func (d *Document) MediaType() string {
	return "document"
}

func (d *Document) MediaFile() *File {
	d.filename = d.FileName
	return &d.File
}

func (d *Document) InputMedia() InputMedia {
	return InputMedia{
		Type:    d.MediaType(),
		Caption: d.Caption,
	}
}
