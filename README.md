# Chatt
> Simple tcp chat in go


### Commands
- `/name <username>` - set username, *anonymous* by default
- `/rooms` - list available rooms
- `/join <name> [password]` - join room
- `/quit` - exit from room, if no room, disconnect
- `/newroom <name> [password] [-h]` - crete new room (-h for hidden)
- `/users` - list users in the room


### TODO
- add secure rooms with encrypted connection

#### Run
`go run cmd/chatt.go -p [port]`


To connect you may use `telnet` or `netcat`