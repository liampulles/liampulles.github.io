package htmlgen

import (
	"bytes"
	"regexp"
	"strconv"
	"unicode"
)

// Good high level ref: https://spec.commonmark.org/current/

// I aim for simplicity and readability here rather then performance. I also don't try to
// cater for the full spec or edge-cases - this only needs to parse my markdown.

type Document struct{}

func ParseMarkdownish(in []byte) (Document, error) {
	// Markdownish consists of line-level block elements and inline elements. I don't allow
	// lazy continuations, which makes parsing simpler.

	// So lets identify the blocks, then identify the inline elements.

	// We can go line-by-line, so let us convert to lines first
	byteLines := bytes.Split(in, []byte{'\n'})
	lines := make([]string, len(byteLines))
	for i, byteLine := range byteLines {
		lines[i] = string(byteLine)
	}

	// Now parse
	parsedLines := parseLines(lines)
	sections := splitSections(parsedLines)
	root := Tree[blockNode]{}
	for _, section := range sections {
		sectionTree := parseBlockTree(section)
		root.children = append(root.children, &sectionTree)
	}

	return Document{}, nil
}

func parseLines(lines []string) []parsedLine {
	allLineBlocks := make([]parsedLine, len(lines))
	for i, l := range lines {
		allLineBlocks[i] = parseLine(l)
	}
	return allLineBlocks
}

type parsedLine struct {
	offset  uint
	markers []marker
	rem     string
}

func parseLine(l string) parsedLine {
	// Pull start-of-line marker
	m, keepPulling, offset := tryPullMarker(l, 0)
	initialOffset := offset
	var markers []marker
	markers = append(markers, m)

	// Pull any remaining markers
	for keepPulling {
		m, keepPulling, offset = tryPullMarker(l, offset)
		if keepPulling {
			markers = append(markers, m)
		}
	}

	return parsedLine{
		offset:  initialOffset,
		markers: markers,
		rem:     l[offset:],
	}
}

func tryPullMarker(line string, offset uint) (marker, bool, uint) {
	cut := splitLine(line[offset:])
	offset += cut
	rem := line[offset:]

	// Different block markers:
	// -> Blank line
	if offset == 0 && rem == "" {
		return blankLineMarker{}, true, uint(len(line))
	}
	// -> Bullet list item
	if elem := bulletListRegex.FindStringSubmatch(rem); len(elem) == 3 {
		return bulletListItemMarker{
			offset: offset,
			bullet: elem[1],
		}, true, offset
	}
	// -> Ordered list item
	if elem := orderedListRegex.FindStringSubmatch(rem); len(elem) == 3 {
		idx, _ := strconv.ParseUint(elem[1], 10, 0)
		return orderedListItemMarker{
			offset: offset + cut,
			idx:    uint(idx),
		}, true, offset
	}
	// -> Code fence
	if elem := codeFenceRegex.FindStringSubmatch(rem); len(elem) == 2 {
		return codeFenceMarker{
			info: elem[1],
		}, true, uint(len(line)) // The rest of the line is chewed up by the info string.
	}

	// Nope, no block elements left.
	return nil, false, offset
}

func lineIsBlank(l string) bool {
	for _, r := range l {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func splitLine(l string) (cut uint) {
	i := uint(0)
	for _, r := range l {
		if !unicode.IsSpace(r) {
			return i
		}
		i++
	}
	// All spaces or empty line
	return i
}

type marker interface{}

type bulletListItemMarker struct {
	offset uint
	bullet string
}

type orderedListItemMarker struct {
	offset uint
	idx    uint
}

type codeFenceMarker struct {
	info string
}

type blankLineMarker struct{}

var bulletListRegex = regexp.MustCompile(`^([\*+-]) (.*)$`)
var orderedListRegex = regexp.MustCompile(`^([0-9]+). (.*)$`)
var codeFenceRegex = regexp.MustCompile(`\x60\x60\x60\s*(.*)`)

// Blank lines close all potential container blocks in Markdownish, so we can split by
// these first to speed things up
func splitSections(parsedLines []parsedLine) [][]parsedLine {
	var sections [][]parsedLine
	var currSection []parsedLine
	for _, pl := range parsedLines {
		if _, isBlankLine := pl.markers[0].(blankLineMarker); isBlankLine && len(currSection) > 0 {
			sections = append(sections, currSection)
			currSection = nil
			continue
		}
		currSection = append(currSection, pl)
	}

	if len(currSection) > 0 {
		sections = append(sections, currSection)
	}
	return sections
}

func parseBlockTree(parsedLines []parsedLine) Tree[blockNode] {
	for _, pl := range parsedLines {
		
	}
}

func 

type blockNode interface{}
