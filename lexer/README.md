# List lexer
Again, this is based on Rob Pike's talk.
Some of the control flow functions will look familiar to anyone who's seen his talk,
because they're often the same functions.

## What if we want to recognize another control symbol?
In Pike's talk, he has to recognize the control symbols `{{` and `}}`.
To do this, he uses 2-lookahead.

Is this necessary? (Answer: No, we could just use extra state functions.)
