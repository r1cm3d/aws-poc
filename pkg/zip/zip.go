package zip

import (
	"archive/zip"
	"aws-poc/internal/protocol"
	"bytes"
	"fmt"
	"strings"
)

type compressor struct{}

func (c compressor) Compress(_ string, files []protocol.File, strToRemove string) ([]byte, error) {
	// TODO: log here
	buffer := &bytes.Buffer{}
	writer := zip.NewWriter(buffer)
	defer writer.Close()

	for _, f := range files {
		fn := strings.ReplaceAll(f.Key, strToRemove, "")
		// TODO: log here
		entry, err := writer.Create(fn)

		if err != nil {
			// TODO: log here
			return nil, err
		}

		// TODO: log
		if _, err := entry.Write(f.Bytes); err != nil {
			// TODO: log here
			return nil, err
		}

		fmt.Printf("%v compressed for %s file", buffer.Bytes(), f.Key)
	}

	return buffer.Bytes(), nil
}
