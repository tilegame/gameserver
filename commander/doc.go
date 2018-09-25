/* 
Package commander enables functions to be called from data structures
that contain the function's name and a list of arguments.


Motivation

The motivation for writing this package was to allow user-accessible
functions to be written more easily.  The main use-case is to call a
function with a JSON message.  The JSON provides the Function's Name,
and the Arguments for that function.  It's very similar to JSON-RPC.

In the context of the game: Let's say another function needs to be
added.  Instead of re-writing a bunch of boilerplate code that
type-checks the incoming JSON messages, the newly written function can
be added to the CommandCenter's Function Map.


How it Works

Reflection is used to analyze the function (which has been written in
Go).  When passed arguments of unknown types (like the JSON message
from a player), the command center handles common error messages such
as "Parameter Type Mismatch" or "Command Not Found".

*/
package commander
