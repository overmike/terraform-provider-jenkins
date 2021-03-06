package jenkins

import (
	"testing"
)

func TestFormatJobName(t *testing.T) {
	inputSimple, inputFolder, inputNested := "job-name", "folder/job-name", "parent/child/job-name"

	// Simple
	actual := formatJobName(inputSimple)
	if actual != inputSimple {
		t.Errorf("Expected %s but received %s", inputSimple, actual)
	}

	// Folder
	actual = formatJobName(inputFolder)
	if actual != "folder/job/job-name" {
		t.Errorf("Expected %s but received %s", inputSimple, actual)
	}

	// Nested
	actual = formatJobName(inputNested)
	if actual != "parent/job/child/job/job-name" {
		t.Errorf("Expected %s but received %s", inputSimple, actual)
	}
}

func TestParseJobName(t *testing.T) {
	inputSimple, inputFolder, inputNested := "job-name", "folder/job-name", "parent/child/job-name"

	// Simple
	actual, actualFolders := parseJobName(inputSimple)
	if actual != inputSimple || len(actualFolders) != 0 {
		t.Errorf("Expected %s with empty folder array but received %s %s", inputSimple, actual, actualFolders)
	}

	// Folder
	actual, actualFolders = parseJobName(inputFolder)
	if actual != inputSimple || len(actualFolders) != 1 || actualFolders[0] != "folder" {
		t.Errorf("Expected %s with single folder array but received %s %s", inputSimple, actual, actualFolders)
	}

	// Nested
	actual, actualFolders = parseJobName(inputNested)
	if actual != inputSimple || len(actualFolders) != 2 || actualFolders[0] != "parent" || actualFolders[1] != "child" {
		t.Errorf("Expected %s with double folder array but received %s %s", inputSimple, actual, actualFolders)
	}
}

func TestTemplateDiff(t *testing.T) {
	// Set up inputs
	inputLeft := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root>Test Case</root>"
	inputRight := "<root>Test {{ .Description }}</root>"

	// Set up Job
	job := resourceJenkinsFolder()
	bag := job.TestResourceData()
	bag.Set("description", "Case")

	if actual := templateDiff("", inputLeft, inputRight, bag); !actual {
		t.Errorf("Expected %s to be considered equal to %s", inputLeft, inputRight)
	}

	// Now try invalid inputs
	inputLeft = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root>Test Incorrect</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); actual {
		t.Errorf("Expected %s to be considered inequal to %s", inputLeft, inputRight)
	}

	inputRight = "<root>Test Incorrect</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); !actual {
		t.Errorf("Expected %s to be considered equal to %s", inputLeft, inputRight)
	}

	inputRight = "<root>Test Even More Incorrect</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); actual {
		t.Errorf("Expected %s to be considered inequal to %s", inputLeft, inputRight)
	}
}
