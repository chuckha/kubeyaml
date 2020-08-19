class Validation {
    constructor(props) {
        this.baseURL = props.baseURL

        this.validate = this.validate.bind(this)
    }

    validate(document, callback) {
        let xhr = new XMLHttpRequest()
        xhr.onreadystatechange = function(e) {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                var status = xhr.status;
                if (status === 0 || (status >= 200 && status < 400)) {
                    callback(xhr.response)
                } else {
                    console.log(xhr.response, status)
                }

            }
        }
        xhr.open("POST", this.baseURL + "/validate", true)
        xhr.send("data="+encodeURI(document))
    }
}

export default Validation
