package filesmodule

import (
	"fmt"
	"strings"
)

type FileId string // 36 (uuid) + 5 (extension) + 1 (dot) = 42

func (ImageId) GormDataType() string { return "varchar(42)" }
func (file *FileId) FileExtension() string {
	if len(string(*file)) == 0 || string(*file)[0] == '.' {
		return ""
	}
	parts := strings.Split(string(*file), ".")
	if len(parts) == 1 {
		return ""
	}

	return parts[len(parts)-1]
}

func (file *FileId) FileName() string {
	if len(string(*file)) == 0 || string(*file)[0] == '.' {
		return string(*file)
	}
	parts := strings.Split(string(*file), ".")
	if len(parts) > 0 {
		fileName, _ := strings.CutSuffix(string(*file), fmt.Sprintf(".%s", parts[len(parts)-1]))
		return fileName
	}
	return string(*file)
}
