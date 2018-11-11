/*
Package commander enables dynamic function calls using strings.


Motivation

The motivation for this package was to make it easier to dynamically
call functions using only a string.

I started by writing insane switch statements based on incoming JSON
messages. This worked well at first, but because tedious to upkeep:
every time I wanted to write a new function, I had to write tons of
boilerplate type checking and validating.

Future experiments involved making a lot of types, and marshalling
JSON messages into those types.
This was well organized and let the encoding/json package deal with
type checking, but it was still tedious for adding new functions.

I look at code generators to turn functions into APIs.
Some of these  become too "magical", and seemed annoying to maintain.
RPC systems seem like overkill.
I just wanted to throw in a string and get the result.

Thus, I wrote this package.
String goes in; Result or Error comes out.




Philosophy

From most important to least important:
	1. Easy to Use.         <- literally the whole point.
	2. Easy to Understand   <- the design is simple.
	3. Easy to Build Over   <- can be used as a foundation.
I want to keep this "magic-free".  It shouldn't feel like tons of
things are happening under the hood.  Basically all it does is do
some type-checking on a string, and call a function with it.




Current Status

Still a work in progress, but I'm quickly started to enjoy this
package a lot.  It's simple, it makes use of meta-programming,
and it doesn't contain tons of "magic".

	1. Type Commander is literally just a map.
	2. Commander.Call() is the only function that you need.
	3. The other methods "CallWith..." are for convenience.



Future Goals

	- support more kinds of functions.
	- allow multiple returns.
	- Support different encodings (not just JSON).


*/
package commander
