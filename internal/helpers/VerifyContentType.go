package helpers

import (
	"io"
	"net/http"
)

func VerifyContentType(file io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)

	n, err := io.ReadFull(file, buffer)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer[:n]), nil
}