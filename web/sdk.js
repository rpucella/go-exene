
// Turn this into a class?

let _ready = false
let _socket = null
let _debug = false

function toggleDebug() {
    for (const elt of document.querySelectorAll("body *")) {
        if (_debug) {
            elt.classList.remove("debug")
        } else {
            elt.classList.add("debug")
        }
    }
    _debug = !(_debug)
}

function connect(url) {
    const socket = new WebSocket(url)
    socket.addEventListener('open', (evt) => {
        console.log('Websocket connection open')
        _ready = true
        initialize()
    })
    socket.addEventListener('close', (evt) => {
        console.log(`Websocket connection closed: ${evt.code} ${evt.reason}`)
        _ready = false
    })
    socket.addEventListener('message', (evt) => {
        const msg = JSON.parse(evt.data)
        if (_debug) {
            console.dir(msg)
        }
        if (msg.type == "message") {
            console.log(`Message: ${msg.text}`)
        } else if (msg.type == "widget") {
            // Clear the document.
            const elt = createWidget(msg.widget)
            while (document.body.firstChild) {
                document.body.firstChild.remove()
            }
            if (elt) {
                document.body.appendChild(elt)
            }
        } else if (msg.type == "update-text") {
            document.querySelector(`#ID${msg.target}`).innerText = msg.text
        } else if (msg.type == "update-size") {
            const elt = document.querySelector(`#ID${msg.target}`)
            elt.style.height = `${msg.height}px`
            elt.style.width = `${msg.width}px`
        } else {
            console.log(`Unknown message: ${evt.data}`)
        }
    })
    _socket = socket
}

function initialize() {
    let timeoutId = null
    // Wait 0.1s without a resize before firing off a resize message.
    const delayResize = 100 
    const height = window.innerHeight
    const width = window.innerWidth
    sendToExene({type: "init", height, width})
    window.addEventListener("resize", (evt) => {
        if (timeoutId) {
            // We already have a timer going, so reset it
            window.clearTimeout(timeoutId)
            timeoutId = null
        }
        timeoutId = window.setTimeout(() => {
            const height = window.innerHeight
            const width = window.innerWidth
            if (_debug) {
                console.log(`Resizing ${width} x ${height}`)
            }
            sendToExene({type: "resize", height, width})
        }, delayResize)
    })
}

function sendToExene(message) {
    if (_ready) {
        _socket.send(JSON.stringify(message))
    }
}

function createWidget(w) {
    const elt = document.createElement(w.tag)
    w.attrs = w.attrs || {}
    w.style = w.style || {}
    w.children = w.children || []
    w.events = w.events || []
    // Defaults.
    elt.style.boxSizing = "border-box"
    if (w.id) {
        elt.setAttribute("id", `ID${w.id}`)
    }
    // Attributes.
    for (const [k, v] of Object.entries(w.attrs)) {
        elt.setAttribute(k, v)
    }
    // Style.
    for (const [k, v] of Object.entries(w.style)) {
        elt.style[k] = v
    }
    // Inner text.
    if (w.text) {
        elt.innerText = w.text
    }
    // Children elements.
    for (const c of w.children) {
        const subelt = createWidget(c)
        elt.appendChild(subelt)
    }
    // Event handlers.
    for (const e of w.events) {
        elt.addEventListener(e, (evt) => {
            sendToExene({type: "event", name: e, target: w.id})
        })
    }
    return elt
}
