class Toast {
    constructor(level, msg) {
        this.level = level
        this.msg = msg
    }

    /**
        Creates new Toast element
        @returns {HTMLButtonElement}
    */
    #create() {
        const t = document.createElement("button")
        t.classList.add("toast", `toast-${this.level}`)
        t.setAttribute("role", "alert")
        t.setAttribute("aria-label", "Close")
        t.addEventListener("click", () => { t.remove() })
        const content = document.createElement("span")
        content.innerText = this.msg
        content.classList.add("toast-message")
        t.appendChild(content)
        return t
    }

    show() {
        const toaster = document.querySelector("#toaster")

        if (!toaster) {
            console.error("Toaster has not been found")
            return
        }
        const t = this.#create()
        toaster.appendChild(t)

        setTimeout(() => {
            t.remove()
        }, 5000)
    }
}

/**
    * @typedef {Object} ToastEvent
    * @property {string} level
    * @property {string} message
*/

/**
    * @param {CustomEvent<ToastEvent>} event 
*/
function handleShowToast(event) {
    if (!event.detail.level || !event.detail.message) {
        console.error("Toast event not valid")
        return
    }

    const toast = new Toast(event.detail.level, event.detail.message)
    toast.show()
}

window.addEventListener("showToast", handleShowToast)
