onload = function() {
    document.querySelector('#main-form').onsubmit = function (event) {
        event.preventDefault()
        const data = new FormData(document.querySelector('#main-form'))
        fetch(
            event.target.action,
            {
                body: JSON.stringify(Object.fromEntries(data)),
                headers: { 'Content-Type': 'application/json' },
                method: event.target.method
            }
        )
        .then(response => response.json())
        .then(json => {
            const id = 1
            var protocol = 'ws://'
            if (window.location.protocol === 'https:') {
                protocol = 'wss://'
            }

            const socket = new WebSocket(`${protocol}${window.location.host}/tasks/${id}/ws`)

            socket.onopen = function (e) {
                console.log(e)
            }
            socket.onerror = function (e) {
                console.error(e)
            }

            socket.onclose = function(e)  {
                console.log("closed", e)
            }

            socket.onmessage = function(e) {
                console.log("message", e)
            }
        })
    }
}
