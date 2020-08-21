import debounce from "debounce";

class Validation {
    constructor(props) {
        this.baseURL = props.baseURL

        this.validate = debounce(this.validate.bind(this), 500)
    }

    validate(document, callback, errorCallback, alwaysCallback) {
        let xhr = new XMLHttpRequest()
        xhr.onreadystatechange = function(e) {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                var status = xhr.status;
                if (status === 0 || (status >= 200 && status < 400)) {
                    callback(xhr.response)
                } else {
                    errorCallback()
                }
                alwaysCallback()
            }
        }
        xhr.open("POST", this.baseURL + "/validate", true)
        xhr.send("data="+encodeURI(document))
    }
}

export default Validation
