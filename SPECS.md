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

### `make_step(stepId, index)`

Example:

```js
socket.on('make_step', function(stepId, index) {

})
```

Server sends to client command to make step. Among parameters there is `id` of `step` to be done and  `index` of pawn with which player is allowed to play.


### `share_step(step)`

Step structure:

```js
{
    step: 0 /* starting from zero */,
    data: "encoded_step" /* step information */
}
```

For example:

```js
socket.emit('share_step', step)
```

Client sends to server step to share with other players.


### `apply_step(step)`

For example:

```js
socket.on('apply_step', function(step) {
     /* do something with step */
})
```

Server sends to client step shared by other peer to apply it on local state.


### `show_error(msg)`

For example:

```js
socket.on("show_error", function(msg) {
    console.log(msg)
})
```

Server sends error msg related to some previous actions or whatever.