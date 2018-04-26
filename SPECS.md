# Specifications and Conventions

## Socket connections

### `create_game(player_name)` event

Example:

```js
socket.emit("create_game", "piaxar")
```

Client sends this event on "Start new game button click". 

`player_name` - user label which will be shown to other users

### `show_endpoint(ip)` event

Example:

```js
socket.on('show_endpoint', function(endpoint) {
     /* do something with ip */
})
```

Servers sends local endpoint address to client to show and share it. 
If ip is null, it means there was problems with establishing local ip address. And it's better for player to check it by his self

### `connect_to_game(endpoint, player_name)`

Example:

```js
socket.emit('connect_to_game', "12.12.12.12:5001", "ivan")
```

Client sends this request if wants to connect to existing game.

### `make_step(index)`


Example:

```js
socket.on('make_step', function(index) {
     /* do something with ip */
})
```

Server sends to client command to make step and `index` of pawn with which player is allowed to play
