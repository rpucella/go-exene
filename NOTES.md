
# Notes

When the websocket disconnects, the SDK should catch and reset the page showing that connection has
been lost.

When the browser refreshes and creates a new socket connection, the old one should be either
dropped, or maybe a different "instance" of the app should be created, and the old one killed?

Alternatives:

- A Widget needs to "attached" to an app via a dispatcher to work
- the dispatch table + the channels are added to the widget by the "attachment" 
- until you call "attach", the widget is "inactive"?

From the original eXene design:

- control channels (control in, control out as a mailbox)

For layout, add the ability to insert + remove child widgets

- this needs to send an update to the webpage to track the changes
- create + insert a new element

Basically, all updates to widgets need to send an update message to the UI to track the change
