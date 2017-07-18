# Experiments with parsing in Go

## Resources
I started by watching Rob Pike's 2011 talk,
[Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE).
I wrote the initial version of my list lexer (`lexer/lexer.go`) while watching his talk.
It was my first stab at a lexer, so that file uses some of Rob's control flow functions.

I'm currently reading through:

- Terence Parr's *Language Implementation Patterns* (which uses Java and ANTLR.)
- *Introduction to Automata Theory, Languages, and Computation* by Hopcroft, Motwani, and Ullman. (3rd).

## Things to cover

- [x] Lexers
- [ ] Parsers
- [ ] Recursive descent (LL(1), LL(k))
- [ ] Abstract syntax trees
