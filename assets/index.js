function setupDropdown() {
    const dropdown = document.querySelector('#app-dropdown')
    const dropdownDetails = dropdown.querySelector('details')
    const selectedAppIcon = document.querySelector('#selected-app .app-icon')
    const selectedAppName = document.querySelector('#selected-app .app-name')
    dropdown.querySelectorAll('.dropdown-item').forEach(item => {
        item.onclick = function (event) {
            const iconUrl = item.dataset.iconUrl
            const name = item.dataset.name
            const id = item.dataset.id

            selectedAppIcon.src = iconUrl
            selectedAppName.innerHTML = name
            dropdown.setAttribute("data-app-id", id)
            dropdownDetails.removeAttribute("open")
        }
    })
}

function enterLoadingState() {
  hideErrorMessage()

  const submitButton = document.querySelector('#submit-button')
  submitButton.innerHTML = '<span>Submitting</span><span class="AnimatedEllipsis"></span>'
  submitButton.setAttribute("aria-disabled", "true")
  submitButton.blur()

  document.querySelector('#app-dropdown details').disabled = true

  document.querySelector('#key-name').disabled = true
}

function exitLoadingState() {
  const submitButton = document.querySelector('#submit-button')
  submitButton.innerHTML = 'Submit'
  submitButton.setAttribute("aria-disabled", "false")

  document.querySelector('#app-dropdown details').disabled = false

  document.querySelector('#key-name').disabled = false
}

function hideErrorMessage() {
  document.querySelector('#error-message').classList.add('hidden')
}

function showErrorMessage(message) {
  const flash = document.querySelector('#error-message')
  flash.innerHTML = message
  flash.classList.remove('hidden')
}

function submitRequest(params) {
  fetch(
      '/tasks',
      {
          body: JSON.stringify(params),
          headers: {
              'Content-Type': 'application/json'
          },
          method: 'POST'
      }
  )
  .then(response => response.json())
  .then(json => {
    if (json['id']) {
      window.location.href = `/tasks/${json['id']}`
    } else {
      showErrorMessage(json['error'] ?? "unknown error")
    }
    exitLoadingState()
  })
}

function setupForm() {
    const showTooltip = function(elementQuery) {
        const element = document.querySelector(elementQuery)
        element.classList.add("tooltipped", "tooltipped-n", "tooltipped-sticky")
    }

    const hideAllTooltip = function() {
        document.querySelectorAll('.tooltipped-sticky').forEach(element => {
            element.setAttribute("class", "")
        })

    }

    const keyNameTextField = document.querySelector('#key-name')
    const appDropdown = document.querySelector('#app-dropdown')

    appDropdown.onfocus = hideAllTooltip
    keyNameTextField.onclick = hideAllTooltip

    document.querySelector('#submit-button').onclick = function () {
        hideAllTooltip()

        const selectedApp = appDropdown.dataset.appId
        if (!selectedApp || selectedApp.length == 0) {
            showTooltip('#selected-app-tooltip')
            return
        }

        const keyName = keyNameTextField.value
        if (!keyName || keyName.length == 0) {
            showTooltip('#key-name-tooltip')
            return
        }

        enterLoadingState()

        const requestParams = {
            app_id: selectedApp,
            key_name: keyName
        }

        submitRequest(requestParams)
    }
}

onload = function() {
    setupDropdown()
    setupForm()
}
