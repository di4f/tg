package tg

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FileType int

const (
	NoFileType FileType = iota
	ImageFileType
)

var (
	UnknownFileTypeErr = errors.New("unknown file type")
)

type File struct {
	path    string
	typ     FileType
	caption string
}

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) withType(typ FileType) *File {
	f.typ = typ
	return f
}

func (f *File) Type() FileType {
	return f.typ
}

func (f *File) Image() *File {
	return f.withType(ImageFileType)
}

func (f *File) Caption(caption string) *File {
	f.caption = caption
	return f
}

func (f *File) NeedsUpload() bool {
	return true
}

func (f *File) UploadData() (string, io.Reader, error) {
	rd, err := os.Open(f.path)
	if err != nil {
		return "", nil, err
	}

	bufRd := bufio.NewReader(rd)

	fileName := filepath.Base(f.path)

	return fileName, bufRd, nil
}

func (f *File) SendData() string {
	return ""
}
func (f *File) Send(
	sid SessionId, bot *Bot,
) (*Message, error) {
	var chattable tgbotapi.Chattable
	cid := sid.ToApi()

	switch f.Type() {
	case ImageFileType:
		photo := tgbotapi.NewPhoto(cid, f)
		photo.Caption = f.caption
		chattable = photo
	default:
		return nil, UnknownFileTypeErr
	}

	msg, err := bot.Api.Send(chattable)

	return &msg, err
}