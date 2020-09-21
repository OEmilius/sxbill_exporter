package controldir

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WalkDirectory(path string) ([]string, error) {
	flist := []string{}
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".bil") {
				flist = append(flist, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
		return flist, err
	}
	return flist, nil
}

//prepare list witch contain files from lastProcessed file and not last file
func GetNotProcessedFiles(lastProcessed string, flist []string) []string {
	newFiles := []string{}
	for i := len(flist) - 2; i >= 0; i-- {
		if flist[i] == lastProcessed {
			newFiles = flist[i+1 : len(flist)-1]
			return newFiles
		}
	}
	//if not found return all files with *.bill extension
	return flist
}
