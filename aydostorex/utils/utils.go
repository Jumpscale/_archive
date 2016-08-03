package utils

import (
	"crypto/md5"
	"fmt"
	"io"

	"os"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("utils")
)

//Hash compute the md5sum of the reader r
func Hash(r io.Reader) (string, error) {
	s, ok := r.(io.Seeker)
	if ok {
		s.Seek(0, os.SEEK_SET)
	}

	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		log.Errorf("Hash, Error reading source: %v", err)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
