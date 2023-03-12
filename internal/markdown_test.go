package internal

import (
	"testing"
)

func TestMarkdownGenerateTitle(t *testing.T) {
	md := MarkdownGenerateTitle("Test", 1)
	if md != "# Test  \n" {
		t.Fatalf("Expected # Test  \n but got %s", md)
	}
}

func TestMarkdownGenerateLink(t *testing.T) {
	md := MarkdownGenerateLink("Test", "https://github.com")
	if md != "[Test](https://github.com)" {
		t.Fatalf("Expected [Test](https://github.com) but got %s", md)
	}
}

func TestMarkdownGenerateListOfStringPointers(t *testing.T) {
	text1 := "Test1"
	text2 := "Test2"
	md := MarkdownGenerateListOfStringPointers([]*string{&text1, &text2})
	if md != "- Test1  \n- Test2  \n" {
		t.Fatalf("Expected - Test1  \n- Test2  \n but got %s", md)
	}
}

func TestMarkdownGenerateList(t *testing.T) {
	md := MarkdownGenerateList([]any{"Test1", "Test2"})
	if md != "- Test1  \n- Test2  \n" {
		t.Fatalf("Expected - Test1  \n- Test2  \n but got %s", md)
	}
}

func TestMarkdownGenerateListItem(t *testing.T) {
	md := MarkdownGenerateListItem("Test")
	if md != "- Test  \n" {
		t.Fatalf("Expected - Test  \n but got %s", md)
	}
}

func TestMarkdownGenerateTableHeader(t *testing.T) {
	md := MarkdownGenerateTableHeader("Test1", "Test2")
	if md != "| Test1 | Test2 |\n| --- | --- |\n" {
		t.Fatalf("Expected | Test1 | Test2 |  \n| --- | --- |  \n but got %s", md)
	}
}

func TestMarkdownGenerateTableRow(t *testing.T) {
	md := MarkdownGenerateTableRow("Test1", "Test2")
	if md != "| Test1 | Test2 |\n" {
		t.Fatalf("Expected | Test1 | Test2 |  \n but got %s", md)
	}
}

func TestMarkdownGenerateTableSeparator(t *testing.T) {
	md := MarkdownGenerateTableSeparator(2)
	if md != "| --- | --- |\n" {
		t.Fatalf("Expected | --- | --- |  \n but got %s", md)
	}
}
