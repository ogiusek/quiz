package filesmodule

import (
	"quizapi/modules/filesmodule"
	"testing"
)

func TestFileExtension(t *testing.T) {
	tests := []struct {
		input    filesmodule.FileId
		expected string
	}{
		{filesmodule.FileId("document.txt"), "txt"},
		{filesmodule.FileId("archive.tar.gz"), "gz"},
		{filesmodule.FileId("image.jpeg"), "jpeg"},
		{filesmodule.FileId("no_extension"), ""},
		{filesmodule.FileId(".hiddenfile"), ""},
		{filesmodule.FileId("anotherfile."), ""},
	}

	for _, test := range tests {
		result := test.input.FileExtension()
		if result != test.expected {
			t.Errorf("For input %s, expected extension %s but got %s", test.input, test.expected, result)
		}
	}
}

func TestFileName(t *testing.T) {
	tests := []struct {
		input    filesmodule.FileId
		expected string
	}{
		{filesmodule.FileId("document.txt"), "document"},
		{filesmodule.FileId("archive.tar.gz"), "archive.tar"},
		{filesmodule.FileId("image.jpeg"), "image"},
		{filesmodule.FileId("no_extension"), "no_extension"},
		{filesmodule.FileId(".hiddenfile"), ".hiddenfile"},
		{filesmodule.FileId("anotherfile."), "anotherfile"},
	}

	for _, test := range tests {
		result := test.input.FileName()
		if result != test.expected {
			t.Errorf("For input %s, expected name %s but got %s", test.input, test.expected, result)
		}
	}
}
