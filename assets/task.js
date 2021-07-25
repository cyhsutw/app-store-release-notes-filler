function timelineEventBackgroundColorClass(event) {
  switch (event.category) {
    case 'start':
      return 'color-bg-backdrop'
    case 'finish':
      return 'color-bg-success'
    case 'info':
      return 'color-bg-info'
    case 'warn':
      return 'color-bg-warning'
    case 'error':
      return 'color-bg-danger'
    default:
      return ''
  }
}

function timelineEventIconSvg(event) {
  switch (event.category) {
    case 'start':
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M8.5.75a.75.75 0 00-1.5 0v5.19L4.391 3.33a.75.75 0 10-1.06 1.061L5.939 7H.75a.75.75 0 000 1.5h5.19l-2.61 2.609a.75.75 0 101.061 1.06L7 9.561v5.189a.75.75 0 001.5 0V9.56l2.609 2.61a.75.75 0 101.06-1.061L9.561 8.5h5.189a.75.75 0 000-1.5H9.56l2.61-2.609a.75.75 0 00-1.061-1.06L8.5 5.939V.75z"></path></svg>`
    case 'finish':
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg>`
    case 'info':
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path fill-rule="evenodd" d="M2 4a1 1 0 100-2 1 1 0 000 2zm3.75-1.5a.75.75 0 000 1.5h8.5a.75.75 0 000-1.5h-8.5zm0 5a.75.75 0 000 1.5h8.5a.75.75 0 000-1.5h-8.5zm0 5a.75.75 0 000 1.5h8.5a.75.75 0 000-1.5h-8.5zM3 8a1 1 0 11-2 0 1 1 0 012 0zm-1 6a1 1 0 100-2 1 1 0 000 2z"></path></svg>`
    case 'warn':
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path fill-rule="evenodd" d="M8.22 1.754a.25.25 0 00-.44 0L1.698 13.132a.25.25 0 00.22.368h12.164a.25.25 0 00.22-.368L8.22 1.754zm-1.763-.707c.659-1.234 2.427-1.234 3.086 0l6.082 11.378A1.75 1.75 0 0114.082 15H1.918a1.75 1.75 0 01-1.543-2.575L6.457 1.047zM9 11a1 1 0 11-2 0 1 1 0 012 0zm-.25-5.25a.75.75 0 00-1.5 0v2.5a.75.75 0 001.5 0v-2.5z"></path></svg>`
    case 'error':
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path fill-rule="evenodd" d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"></path></svg>`
    default:
      return ''
  }
}


function createTimelineElement(event) {
  const colorClass = timelineEventBackgroundColorClass(event)
  const iconSvg = timelineEventIconSvg(event)
  const message = event.message
  const html = `
    <div class="TimelineItem closed-timeline-item">
      <div class="TimelineItem-badge ${colorClass} color-icon-secondary">
        ${iconSvg}
      </div>
    <div class="TimelineItem-body multiline">${message}</div>
  `

  const template = document.createElement('template')
  template.innerHTML = html.trim()
  return template.content.firstChild
}

function setupWebsocket() {
  const container = document.querySelector('#events')
  const loadingIndicator = document.querySelector('#loading')

  const taskId = container.dataset.taskId
  var protocol = "wss://"
  if (window.location.protocol === 'http:') {
    protocol = "ws://"
  }

  const socket = new WebSocket(`${protocol}//${window.location.host}/tasks/${taskId}/ws`)

  socket.onmessage = function (event) {
    const logEvent = JSON.parse(event.data)
    const element = createTimelineElement(logEvent)
    container.insertBefore(element, loadingIndicator)
    scrollTo(0, document.body.scrollHeight)
  }

  socket.onclose = function () {
    container.removeChild(loadingIndicator)
  }
}

onload = function () {
  setupWebsocket()
}
