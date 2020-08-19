import parseCST from 'yaml/parse-cst'
import safeHtml from 'safe-html'

const config = {
    allowedTags: [],
    allowedAttributes: {}
}


// Some kinda weird mutation going on here.
class Document {
    // ErrorObj: Error, Key
    constructor(doc, errorObj) {
        this.doc = parseCST(doc)
        // will only be one document
        let items = this.doc[0].contents[0]

        this.findNode = this.findNode.bind(this)
        this.findNode(errorObj["Key"], items)
    }

    // findNode parses the YAML
    findNode(key, elements) {
        if (elements === undefined || elements === null) {
            return
        }
        console.log("Key: ", key, "elements:", elements, elements.type)
        let parsedKeys = key.split(".")
        if (elements.type === "PLAIN") {
            elements.value = `<span class="has-text-danger">${safeHtml(elements.strValue, config)}</span>`
            return
        }
        if (elements.type === "MAP") {
            let items = elements.items
            for (let i=0; i < items.length; i++) {
                if (items[i].strValue === parsedKeys[0]) {
                    items[i].value = `<span class="has-text-danger">${safeHtml(items[i].strValue, config)}</span>`
                    if (items[i+1].type === "MAP_VALUE") {
                        this.findNode(parsedKeys.slice(1).join("."), items[i+1].node)
                    }
                    console.log("unexpected type", items[i].type)
                    return
                }
            }
        }
    }
}

export default Document
