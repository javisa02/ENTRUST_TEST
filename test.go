package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	MaxCharsPerLine = 80
	MaxLinesPerPage = 25
)

/*
 * Document represents a document to be paginated
 */
type Document struct {
	Pages []*Page
}

/*
 * Page represents a single page in the document
 */
type Page struct {
	Number int
	Lines  []string
}

/*
 * AddPage adds a page to the document
 */
func (doc *Document) AddPage(page *Page) {
	doc.Pages = append(doc.Pages, page)
}

/*
 * SaveToFile saves the document to a file
 */
func (doc *Document) SaveToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	//For each page in the document
	for _, page := range doc.Pages {
		fmt.Fprintf(writer, "--- Page %d --- \n", page.Number) // Write the page number to the file

		// For each line in the page
		for _, line := range page.Lines {
			fmt.Fprintln(writer, line)
		}
		fmt.Fprintln(writer, "")
	}

	writer.Flush() //Write the buffered data to the file
	return nil
}

/*
 * paginateDocument paginates a document
 */
func paginateDocument(fileName string) (*Document, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := &Document{}
	page := &Page{Number: 1}
	lineCount := 0
	reader := bufio.NewScanner(file)

	for reader.Scan() {
		line := reader.Text()
		line = strings.TrimRight(line, "\n")

		words := strings.Fields(line)
		currentLine := ""

		for _, word := range words {
			wordLength := len(word)

			// Check if adding the word exceeds the maximum line length, taking into account a blank space
			if len(currentLine)+wordLength+1 > MaxCharsPerLine {
				// Check if the next word can fit in the current line with a blank space
				if len(currentLine)+1 < MaxCharsPerLine {
					currentLine = currentLine + " " + word
				} else {
					// If the current line is not empty, add it to the page
					if currentLine != "" {
						page.Lines = append(page.Lines, currentLine)
						lineCount++
						currentLine = word
					} else {
						// If the current line is empty, start a new line with the word
						currentLine = word
					}
				}
			} else {
				// Add the word to the current line if its length is less or equal to the maximum line length
				if currentLine == "" {
					currentLine = word
				} else {
					currentLine += " " + word
				}

			}

			// Check if the current line is completely filled
			if len(currentLine) == MaxCharsPerLine {
				page.Lines = append(page.Lines, currentLine)
				lineCount++
				currentLine = ""
			}

			// Check if the maximum lines per page limit has been reached
			if lineCount == MaxLinesPerPage {
				doc.AddPage(page)
				page = &Page{Number: len(doc.Pages) + 1}
				lineCount = 0
			}

		}

	}

	// Add the last page to the document if it's not empty
	if len(page.Lines) > 0 {
		doc.AddPage(page)
	}

	return doc, nil
}

/*
 * main is the entry point of the program
 */
func main() {
	fileName := "document.txt"
	doc, err := paginateDocument(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = doc.SaveToFile("result.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Paginated document has been saved to result.txt")
}
