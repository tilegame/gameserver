/*
package cookiez uses gorilla secure cookies to store a player's ID.

The motivation for using cookies is to keep a player logged in even if their
websocket connection changes.  This noticeably happens on mobile devices,
which often disconnect when the browser page in no longer in the foreground.

How it Works

When a player first logs in, they are assigned a player ID.
The id is randomly generated and saved on the server.
This is not enough by itself, because you could just guess different id's
until one works, and then use that to control other players.

To solve this problem, there can be a cookie saved by the client, which
can remain secure assuming the connection is already secured through TLS.
When a player issues a command, the message will contain both a playerid,
and a validation.

If the playerid and the validation match, then it is assumed that the
command originated from the same client that logged in.
*/
package cookiez
