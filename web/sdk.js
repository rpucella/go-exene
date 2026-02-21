
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
        const elt = document.createElement("div")
        elt.style.position = "absolute"
        elt.style.left = "0"
        elt.style.right = "0"
        elt.style.top = "0"
        elt.style.bottom = "0"
        elt.style.zIndex = "100"
        elt.style.opacity = "0.8"
        elt.style.padding = "100px"
        elt.style.textAlign = "center"
        elt.style.fontSize = "48px"
        elt.style.color = "white"
        elt.style.backgroundColor = "black"
        elt.innerText = "Connection lost"
        document.body.appendChild(elt)
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
            document.querySelector(`#widget-${msg.target}`).innerText = msg.text
        } else if (msg.type == "update-size") {
            const elt = document.querySelector(`#widget-${msg.target}`)
            elt.style.height = `${msg.height}px`
            elt.style.width = `${msg.width}px`
        } else if (msg.type == "insert-child") {
            const elt = document.querySelector(`#widget-${msg.target}`)
            const newElt = createWidget(msg.widget)
            elt.insertBefore(newElt, elt.childNodes[msg.index])
        } else if (msg.type == "append-child") {
            const elt = document.querySelector(`#widget-${msg.target}`)
            const newElt = createWidget(msg.widget)
            elt.appendChild(newElt)
        } else if (msg.type == "delete-child") {
            const elt = document.querySelector(`#widget-${msg.target}`)
            elt.removeChild(elt.childNodes[msg.index])
        } else if (msg.type == "hide-child") {
            const elt = document.querySelector(`#widget-${msg.target}`).childNodes[msg.index]
            if (elt.style.display !== "none") {
                elt.setAttribute("data-display-save", elt.style.display)
                elt.style.display = "none"
            }
        } else if (msg.type == "unhide-child") {
            const elt = document.querySelector(`#widget-${msg.target}`).childNodes[msg.index]
            if (elt.style.display === "none" && elt.hasAttribute("data-display-save")) {
                elt.style.display = elt.getAttribute("data-display-save")
                elt.removeAttribute("data-display-save")
            }
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
    if (!w.tag) {
        // Skip empty widgets.
        // (They come from doubly created realizers!)
        return
    }
    const elt = document.createElement(w.tag)
    w.attrs = w.attrs || {}
    w.style = w.style || {}
    w.children = w.children || []
    w.events = w.events || []
    // Defaults.
    elt.style.boxSizing = "border-box"
    if (w.id) {
        elt.setAttribute("id", `widget-${w.id}`)
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
