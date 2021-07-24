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
        .then(response => { console.log(response.json() )})
    }
}
