package core

import (
	"github.com/zhangyiming748/finder"
)

func BatchChangeTime(dirPath string, times *FileTimes) error {
	files := finder.FindAllFiles(dirPath)
	for _, file := range files {
		err := SetFileTimes(file, times)
		if err != nil {
			return err
		}
	}
	return nil
}
