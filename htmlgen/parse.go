package htmlgen

import (
	"bytes"
	"regexp"
	"strconv"
	"unicode"
)

// Good high level ref: https://spec.commonmark.org/0.30/#appendix-a-parsing-strategy

// I aim for simplicity and readability here rather then performance. I also don't try to
// cater for the full spec or edge-cases - this only needs to parse my markdown.

type Document struct{}

func ParseMarkdownish(in []byte) (Document, error) {
	// Roughly, a markdown file consists of blocks (at the highest level, separated by a blank lines).
	// Each block consists of child blocks recursively, and ultimately some text.
	// The text then consists of inline elements.

	// So lets identify the blocks, then identify the inline elements.

	// We can go line-by-line, so let us convert to lines first
	byteLines := bytes.Split(in, []byte{'\n'})
	lines := make([]line, len(byteLines))
	for i, byteLine := range byteLines {
		lines[i] = line(byteLine)
	}

	// Parse the blocks
	root := parseLinesToBlockTree(lines)

	return Document{}, nil
}

func parseLinesToBlockTree(lines []line) *block {
	// We can do this in one pass of the lines (though many passes of the blocks).
	// Block identifiers are always at the front of the line, which will help.
	//
	// First we "consume continuation markers" from the front of the line. E.g.
	// if we are currently adding on to a list inside of a block-quote, we'll chomp
	// of a leading > and -. We do this by matching markers against our tree in
	// a breadth first order. We keep track of the last one matched.
	//
	// If there are any special markers left, then it indicates the need for a
	// new block. We add it as a child of the last matched block identified above,
	// and move the parent's text over to it. We also close any open blocks in the tree.
	root := &Tree[block]{
		item: block{typ: document},
	}
	for _, l := range lines {
		// Check to close first, then check to open.
		// We need to keep track of the last one closed (if any)
		// because that is where new blocks will be added.
		lastMatched := consumeContinuationMarkers(root, l)
	}
}

func consumeContinuationMarkers(t *Tree[block], l line) (lastMatched *Tree[block]) {
	indent, rem := splitLine(l)
	t.IterateBreadthFirst(func(blockNode *Tree[block]) (stop bool) {
		// Only interested in open blocks
		if blockNode.item.closed {
			return false
		}

		//
	})
}

type line string

func lineIsBlank(l line) bool {
	for _, r := range l {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func splitLine(l line) (indent string, rem string) {
	var indentRunes []rune
	var remRunes []rune
	splitPointReached := false
	for _, r := range l {
		if !splitPointReached && unicode.IsSpace(r) {
			indentRunes = append(indentRunes, r)
			continue
		}
		splitPointReached = true
		remRunes = append(remRunes, r)
	}
	return string(indentRunes), string(remRunes)
}

func blockShouldRemainOpen(typ blockType, rem string) bool {

}

type block[T any] struct {
	blockType
	closed bool
}

type blockType interface {
	matchesLine(rem string) bool
}

type document struct{}

func (document) matchesLine(rem string) bool {
	// Document stays open as long as there are lines
	return true
}

type paragraph []string

func (paragraph) matchesLine(rem string) bool {
	// Must be some text
	return len(rem) > 0
}

type list struct {
	isOrdered bool
	bullet    string
}

func (l list) matchesLine(rem string) bool {
	// Only matches if it is a list item of the same type
	listMarker, _, isList := splitListMarker(rem)
	// -> Not a list marker at all
	if !isList {
		return false
	}
	// -> Also an ordered list
	if l.isOrdered && listMarker.isOrdered {
		return true
	}
	// -> Matching bullet
	if l.bullet == listMarker.bullet {
		return true
	}
	// Ok, not matching lists
	return false
}

type listItem struct{}

func (listItem) matchesLine(rem string) bool {
	// Never matches
	return false
}

type blockQuote struct{}

func (blockQuote) matchesLine(rem string) bool {
	// Must be some text
	return len(rem) > 0
}

type code struct{}

func (code) matchesLine(rem string) bool {
	// If we hit a code fence marker, then se should close
	_, isCode := splitFencedCodeMarker(rem)
	if isCode {
		return false
	}
	return true
}

type listMarker struct {
	isOrdered bool
	bullet    string
	idx       uint
}

func splitListMarker(rem string) (listMarker, string, bool) {
	// Try bullets
	if bulletElem := bulletListMarkerRegex.FindStringSubmatch(rem); len(bulletElem) >= 3 {
		bullet, rest := bulletElem[1], bulletElem[2]
		return listMarker{
			bullet: bullet,
		}, rest, true
	}

	// Try ordered
	if orderedElem := orderedListMarkerRegex.FindStringSubmatch(rem); len(orderedElem) >= 3 {
		idxStr, rest := orderedElem[1], orderedElem[2]
		idx, _ := strconv.ParseUint(idxStr, 10, 0)
		return listMarker{
			isOrdered: true,
			idx:       uint(idx),
		}, rest, true
	}

	// Ok, its not a list item
	return listMarker{}, rem, false
}

func splitFencedCodeMarker(rem string) (string, bool) {
	if infoElem := codeFenceRegex.FindStringSubmatch(rem); len(infoElem) >= 2 {
		info := infoElem[1]
		return info, true
	}

	return rem, false
}

var bulletListMarkerRegex = regexp.MustCompile(`^([\*+-]) (.*)$`)
var orderedListMarkerRegex = regexp.MustCompile(`^([0-9]+). (.*)$`)
var codeFenceRegex = regexp.MustCompile(`\x60\x60\x60\s*(.*)`)
