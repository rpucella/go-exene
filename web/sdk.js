
function printString(s) {
    const elt = document.createElement("p")
    elt.innerText = s
    document.body.appendChild(elt)
}
    
// Turn this into a class?

let _ready = false
let _socket = null

function connect(url) {
    const socket = new WebSocket(url)
    socket.addEventListener('open', (evt) => {
        console.log('Websocket connection open')
        _ready = true
    })
    socket.addEventListener('close', (evt) => {
        console.log(`Websocket connection closed: ${evt.code} ${evt.reason}`)
        _ready = false
    })
    socket.addEventListener('message', (evt) => {
        const msg = JSON.parse(evt.data)
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
        } else if (msg.type == "update") {
            if (msg.text) {
                document.querySelector(`#ID${msg.target}`).innerText = msg.text
            } else {
                console.log(`Unknown update: ${evt.data}`)
            }
        } else {
            console.log(`Unknown message: ${evt.data}`)
        }
    })
    _socket = socket
}

function sendToExene(message) {
    if (_ready) {
        _socket.send(JSON.stringify(message))
    }
}

function createWidget(w) {
    let elt
    switch (w.type) {
    case "button":
        elt = document.createElement("button")
        elt.setAttribute("id", `ID${w.id}`)
        elt.style.boxSizing = "border-box"
        for (const entry of Object.entries(w.style)) {
            elt.style[entry[0]] = entry[1]
        }
        elt.innerText = w.label
        elt.addEventListener("click", (evt) => {
            msg = {
                "type": "event",
                "event": "click",
                "target": w.id
            }
            _socket.send(JSON.stringify(msg))
        })
        return elt
        
    case "text":
        elt = document.createElement("span")
        elt.setAttribute("id", `ID${w.id}`)
        elt.style.boxSizing = "border-box"
        for (const entry of Object.entries(w.style)) {
            elt.style[entry[0]] = entry[1]
        }
        elt.innerText = w.text
        return elt
        
    case "gap":
        elt = document.createElement("div")
        elt.setAttribute("id", `ID${w.id}`)
        elt.style.boxSizing = "border-box"
        elt.style.height = w.size
        elt.style.width = w.size
        for (const entry of Object.entries(w.style)) {
            elt.style[entry[0]] = entry[1]
        }
        return elt
        
    case "layout":
        elt = document.createElement("div")
        elt.setAttribute("id", `ID${w.id}`)
        elt.style.boxSizing = "border-box"
        for (const entry of Object.entries(w.style)) {
            elt.style[entry[0]] = entry[1]
        }
        elt.style.display = "flex"
        elt.style.flexDirection = w.direction.startsWith("column") ? "column" : "row"
        if (w.direction.endsWith("top") || w.direction.endsWith("left")) {
            elt.style.alignItems = "flex-start"
        } else if (w.direction.endsWith("bottom") || w.direction.endsWith("right")) {
            elt.style.alignItems = "flex-end"
        } else {
            elt.style.alignItems = "center"
        }
        for (const w2 of w.widgets) {
            const elt2 = createWidget(w2)
            elt.appendChild(elt2)
        }
        return elt

    default:
        console.log(`Unknown widget type ${w.type} in ${JSON.stringify(w)}`)
        return null
    }
}
