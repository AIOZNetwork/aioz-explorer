package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func ConvertNanoTimestampToMilliSecond(nano int64) int64 {
	return nano / 1E6
}

func ResetNodeDB(nodeDir string) error {
	dir := path.Join(nodeDir, "data")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, f := range files {
		if !strings.Contains(f.Name(), ".json") {
			if err = os.RemoveAll(path.Join(nodeDir, "data", f.Name())); err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}

func JSONifyObject(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Min2Int(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
