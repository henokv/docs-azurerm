package internal

import (
	"fmt"
	"os"
	"strings"
)

type Markdown struct {
	docsDir string
}

func NewMarkdown() *Markdown {
	return &Markdown{}
}

func (client *Markdown) writeToFile(content string, destination string) (err error) {
	path := destination[:strings.LastIndex(destination, "/")]
	os.MkdirAll(path, 0770)
	file, err := os.Create(destination)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		file, err = os.Open(destination)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (client *Markdown) GenerateTitle(title string, level int) string {
	return fmt.Sprintf("%s %s  \n", strings.Repeat("#", level), title)
}

func (client *Markdown) GenerateLink(title string, link string) string {
	return fmt.Sprintf("[%s](%s)", title, strings.ReplaceAll(link, " ", "%20"))
}

func (client *Markdown) GenerateListOfStringPointers(items []*string) (listText string) {
	for _, item := range items {
		listText += client.GenerateListItem(*item)
	}
	return listText
}

//func MarkdownGenerateListOfPointers(items []*any) (listText string) {
//	for _, item := range items {
//		listText += GenerateListItem(*item)
//	}
//	return listText
//}

func (client *Markdown) GenerateList(items []any) (listText string) {
	for _, item := range items {
		listText += client.GenerateListItem(item)
	}
	return listText
}

func (client *Markdown) GenerateListItem(itemText any) string {
	return fmt.Sprintf("- %s  \n", itemText)
}

//func MarkdownGenerateTable(itemText any) string {
//	return fmt.Sprintf("- %s  \n", itemText)
//}

func (client *Markdown) GenerateTableHeader(headers ...string) string {
	markdown := client.GenerateTableRow(headers...)
	separator := client.GenerateTableSeparator(len(headers))
	return markdown + separator
}

func (client *Markdown) GenerateTableSeparator(columnCount int) string {
	markdown := "|"
	for i := 0; i < columnCount; i++ {
		markdown += " --- |"
	}
	markdown += "\n"
	return markdown
}

func (client *Markdown) GenerateTableRow(fields ...string) string {
	markdown := "|"
	for _, field := range fields {
		markdown += fmt.Sprintf(" %s |", field)
	}
	markdown += "\n"
	return markdown
}
