// Written while watching Rob Pike's 2011 talk, Lexical Scanning in Go
// https://www.youtube.com/watch?v=HxaD_trXwRE
// Some functions taken directly from the talk.

package lexer

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

type itemType int

const eof = -1

// Basically an autoincrementing enum
const (
    itemError itemType = iota
    itemLeftBracket
    itemRightBracket
    itemComma
    itemSymbol
    itemEOF
)

type item struct {
    typ itemType
    val string
}

// Pretty printing for debugging
func (i item) String() string {
    switch i.typ {
    case itemEOF:
        return "EOF"
    case itemError:
        return i.val
    }
    if len(i.val) > 10 {
        // Truncate to 10 characters
        return fmt.Sprintf("%.10q...", i.val)
    }
    return fmt.Sprintf("%q", i.val)
}

// Current state is a function that returns another state function
type stateFn func(*lexer) stateFn

// Ripped straight from Pike's talk
type lexer struct {
    name string     // used only for error reports.
    input string    // the string being scanned.
    start int       // start position of this item.
    pos int         // current position in the input.
    width int       // width of last rune read. Needed for UTF-8.
    items chan item // channel of scanned item.
}

//TODO this changes
func Lex(name, input string) (*lexer, chan item) {
    // A lexer initializes itself to a string,
    l := &lexer{
        name: name,
        input: input,
        items: make(chan item),
    }
    // then launches the state machine as a goroutine,
    go l.run()
    // returning the lexer itself and a channel of items.
    return l, l.items
}

// State machine function; runs as a goroutine
func (l *lexer) run() {
    // Remember, states are functions - so state & lexText are functions
    for state := lexText; state != nil; {
        state = state(l)
    }
    // Close the items channel, as no more tokens will be delivered.
    close(l.items)
}

// Pass an item back to client on the provided channel
func (l *lexer) emit(t itemType) {
    l.items <- item{t, l.input[l.start:l.pos]}
    l.start = l.pos
}

// This is our entry point.
func lexText(l *lexer) stateFn {
    // Expect a left bracket or EOF
    switch l.next() {
    case '[':
        l.emit(itemLeftBracket)
        return lexInsideList
    case eof:
        l.emit(itemEOF) // It's useful to provide EOF as a token.
        return nil      // Stop the run loop.
    default:
        return l.errorf("unexpected initial symbol")
    }
}

/**
 * Flow control
 */

// Consume next UTF-8 rune
func (l *lexer) next() rune {
    if l.pos >= len(l.input) {
        l.width = 0
        return eof
    }
    runeVal, width := utf8.DecodeRuneInString(l.input[l.pos:])
    // Can't seem to use this in the assignment; compiler screaming
    l.width = width
    l.pos += l.width
    return runeVal
}

func (l *lexer) ignore() {
    l.start = l.pos
}

// Step back one rune.
// Can't call more than once per next(), since runes vary in width.
func (l *lexer) backup() {
    l.pos -= l.width
}

// Get next rune without consuming it.
func (l *lexer) peek() rune {
    runeVal := l.next()
    l.backup()
    return runeVal
}

// "accept" in Pike's talk; I'm adapting for LIP
// Accept the next rune if it's in the valid set.
func (l *lexer) consume(valid string) bool {
    if strings.IndexRune(valid, l.next()) >= 0 {
        return true
    }
    l.backup()
    return false
}

// Consume until failure, then backup for failed rune
func (l *lexer) consumeMany(valid string) {
    for strings.IndexRune(valid, l.next()) >= 0 {
    } // Do nothing but call .next()
    l.backup()
}

func (l *lexer) consumeWS() {
    l.consumeMany(" \t\r\n")
    l.ignore()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
    l.items <- item{
        itemError,
        fmt.Sprintf(format, args...),
    }
    return nil
}

func isSymbolRune(r rune) bool {
    runes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    return 0 <= strings.IndexRune(runes, r)
}

/**
 * Accept functions
 */

// Consumes characters, then accept comma or ]
func lexSymbol(l *lexer) stateFn {
    chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    l.consumeMany(chars)
    if l.pos > l.start {
        // Only emit if consumeMany found at least one rune
        // (This should always be the case, though, given state machine model)
        l.emit(itemSymbol)
    }

    // Ignore white space before comma or bracket
    l.consumeWS()
    switch l.next() {
    case ',':
        l.emit(itemComma)
        return lexInsideList
    case ']':
        l.emit(itemRightBracket)
        return lexAfterList
    default:
        return l.errorf("unexpected rune")
    }
}

func lexInsideList(l *lexer) stateFn {
    l.consumeWS()
    switch r := l.next(); {
    case isSymbolRune(r):
        // Using our lookahead here
        l.backup()
        return lexSymbol
    case r == '[':
        l.emit(itemLeftBracket)
        return lexInsideList
    case r == ']':
        l.emit(itemRightBracket)
        return lexAfterList
    case r == eof:
        return l.errorf("unclosed list")
    default:
        return l.errorf("unexpected rune")
    }
}

func lexAfterList(l *lexer) stateFn {
    // Just got ]. Now, expect: WS, comma, ], or EOF
    l.consumeWS()
    switch l.next() {
    case ',':
        l.emit(itemComma)
        return lexInsideList
    case ']':
        l.emit(itemRightBracket)
        return lexAfterList
    case eof:
        l.emit(itemEOF)
        return nil
    default:
        return l.errorf("unexpected rune")
    }
}
