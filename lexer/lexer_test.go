package lexer

import (
    "fmt"
    "testing"
)

// Accept a single item from a channel and test that it is of itemType typ
func assertItemType(c <-chan item, typ itemType, t *testing.T) {
    i := <-c
    if i.typ != typ {
        fmt.Println("Wrong item type")
        t.Fail()
    }
}

// Trivial test; I wrote this one to get a feel for the testing package
func TestIsSymbolRune(t *testing.T) {
    if ! isSymbolRune('a') {
        t.Fail()
    }
    if isSymbolRune('\n') {
        t.Fail()
    }
}

func TestLexSymbol(t *testing.T) {
    l := &lexer{
        input: "asdf",
        items: make(chan item),
    }
    go lexSymbol(l)
    item1 := <-l.items
    if item1.typ != itemSymbol {
        fmt.Println("Expected itemSymbol, got:", item1)
        t.Fail()
    }
    if item1.val != "asdf" {
        fmt.Println("Expected value 'asdf', got:", item1)
        t.Fail()
    }
    item2 := <-l.items
    if item2.typ != itemError {
        fmt.Println("Expected itemError, got:", item2)
        t.Fail()
    }
}

func TestLexEmptyList(t *testing.T) {
    _, items := Lex("lex empty list", "[]")
    item1 := <-items
    if item1.typ != itemLeftBracket {
        fmt.Println("Expected itemLeftBracket, got:", item1)
        t.Fail()
    }
    if item1.val != "[" {
        fmt.Println("Expected value '[', got:", item1)
        t.Fail()
    }
    item2 := <-items
    if item2.typ != itemRightBracket {
        fmt.Println("Expected itemRightBracket, got:", item2)
        t.Fail()
    }
    item3 := <-items
    if item3.typ != itemEOF {
        fmt.Println("Expected itemEOF, got:", item3)
        t.Fail()
    }
    //TODO Test channel closed
}

func TestLexUnitList(t *testing.T) {
    _, items := Lex("lex unit list", "[ asdf ]")
    assertItemType(items, itemLeftBracket, t)
    assertItemType(items, itemSymbol, t)
    assertItemType(items, itemRightBracket, t)
    assertItemType(items, itemEOF, t)
    //TODO Test channel closed
}

func TestLexTriplet(t *testing.T) {
    _, items := Lex("lex unit list", "[ a,b, c ]")
    assertItemType(items, itemLeftBracket, t)
    assertItemType(items, itemSymbol, t)
    assertItemType(items, itemComma, t)
    assertItemType(items, itemSymbol, t)
    assertItemType(items, itemComma, t)
    assertItemType(items, itemSymbol, t)
    assertItemType(items, itemRightBracket, t)
    assertItemType(items, itemEOF, t)
}

func TestLexNestedList(t *testing.T) {
    _, items := Lex("lex nested list", "[ [], a ]")
    assertItemType(items, itemLeftBracket, t)
    assertItemType(items, itemLeftBracket, t)
    assertItemType(items, itemRightBracket, t)
    assertItemType(items, itemComma, t)
    assertItemType(items, itemSymbol, t)
    assertItemType(items, itemRightBracket, t)
    assertItemType(items, itemEOF, t)
}

func TestBrokenListsFail(t *testing.T) {
    _, items := Lex("lex right open list", "[ a")
    assertItemType(items, itemLeftBracket, t)
    assertItemType(items, itemSymbol, t)
    assertItemType(items, itemError, t)
    _, items = Lex("lex unstarted list", "a ]")
    assertItemType(items, itemError, t)
}
