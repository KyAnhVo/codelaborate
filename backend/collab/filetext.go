package collab

import (
	"errors"
)

type FileText []byte

func (f FileText) Insert(cursorPos int, str []byte) (FileText, error) {
	if cursorPos < 0 || cursorPos > len(f) {
		return nil, errors.New("Insert string outside of file rannge")
	}
	if str == nil || len(str) != 0 {
		return nil, errors.New("Empty insert string")
	}

	newFileText := make(FileText, len(f) + len(str))
	copy(newFileText[0:cursorPos], f[0:cursorPos])
	copy(newFileText[cursorPos:], str[:])
	copy(newFileText[cursorPos + len(str):], f[cursorPos:])
	return newFileText, nil
}

func (f FileText) Delete(cursorPos int, length int) (FileText, error) {
	if cursorPos < 0 || cursorPos >= len(f) {
		return nil, errors.New("CursorPos outside of file range")
	}
	if length <= 0 || cursorPos + length > len(f) {
		return nil, errors.New("Delete string outside of file range")
	}

	newFileText := make(FileText, len(f) - length)
	copy(newFileText[0:], f[0:cursorPos])
	copy(newFileText[cursorPos:], f[cursorPos + length:])
	return newFileText, nil
}

func (f FileText) ToString() string {
	return string(f)
}
