package internal

import (
	"fmt"
	"strings"
)

func MarkdownGenerateTitle(title string, level int) string {
	return fmt.Sprintf("%s %s  \n", strings.Repeat("#", level), title)
}

func MarkdownGenerateLink(title string, link string) string {
	return fmt.Sprintf("[%s](%s)", title, strings.ReplaceAll(link, " ", "%20"))
}

func MarkdownGenerateListOfStringPointers(items []*string) (listText string) {
	for _, item := range items {
		listText += MarkdownGenerateListItem(*item)
	}
	return listText
}

//func MarkdownGenerateListOfPointers(items []*any) (listText string) {
//	for _, item := range items {
//		listText += MarkdownGenerateListItem(*item)
//	}
//	return listText
//}

func MarkdownGenerateList(items []any) (listText string) {
	for _, item := range items {
		listText += MarkdownGenerateListItem(item)
	}
	return listText
}

func MarkdownGenerateListItem(itemText any) string {
	return fmt.Sprintf("- %s  \n", itemText)
}

//func MarkdownGenerateTable(itemText any) string {
//	return fmt.Sprintf("- %s  \n", itemText)
//}

func MarkdownGenerateTableHeader(headers ...string) string {
	markdown := MarkdownGenerateTableRow(headers...)
	separator := MarkdownGenerateTableSeparator(len(headers))
	return markdown + separator
}

func MarkdownGenerateTableSeparator(columnCount int) string {
	markdown := "|"
	for i := 0; i < columnCount; i++ {
		markdown += " --- |"
	}
	markdown += "\n"
	return markdown
}

func MarkdownGenerateTableRow(fields ...string) string {
	markdown := "|"
	for _, field := range fields {
		markdown += fmt.Sprintf(" %s |", field)
	}
	markdown += "\n"
	return markdown
}
