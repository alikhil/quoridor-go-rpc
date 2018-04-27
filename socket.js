function init(socket) {
    socket.on("show_endpoint", console.log)
    socket.on("make_step", console.log)
    socket.on("apply_step", console.log)
    socket.on("show_error", console.error)
}

socket = io()

init(socket)