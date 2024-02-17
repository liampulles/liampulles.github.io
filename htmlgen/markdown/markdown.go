package markdown

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

// Parse markdown. Don't go wild - ideally want to use our DSL to build main elements, just use basic markdown.

var md = goldmark.New()

func Parse(s string) (template.HTML, error) {
	var sb strings.Builder
	err := md.Convert([]byte(s), &sb)
	if err != nil {
		return "", err
	}
	return template.HTML(sb.String()), nil
}

// Split lines, then group those lines into blocks.
// Blocks are separated by empty lines.
func splitIntoLineBlocks(md string) [][]string {
	lines := strings.Split(md, "\n")
	var blocks [][]string
	var currBlock []string
	for _, line := range lines {
		if isBlank(line) && len(currBlock) > 0 {
			blocks = append(blocks, currBlock)
			currBlock = nil
		} else {
			currBlock = append(currBlock, line)
		}
	}
	if len(currBlock) > 0 {
		blocks = append(blocks, currBlock)
	}
	return blocks
}

func isBlank(s string) bool {
	s = strings.TrimSpace(s)
	return s == ""
}

var (
	codeFenceRegex = regexp.MustCompile(`^\s*\x60\x60\x60\s*(\w*)\s*$`)
)

func parseBlock(lines []string) (Block, error) {
	// In my simple model, I determine a block by matching the first line to some regex.
	return nil, nil
}

type Document struct {
	Blocks []Block
}

type Block interface{}
