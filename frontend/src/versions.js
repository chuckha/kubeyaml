

// TODO: rewrite as promise
function fetchVersions(baseURL, success, error, always) {
    let xhr = new XMLHttpRequest()
    xhr.onreadystatechange = function(e) {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            var status = xhr.status
            if (status === 0 || (status >= 200 && status < 400)) {
                success(xhr.response)
            } else {
                error()
            }
            always()
        }
    }
    xhr.open("GET", baseURL + "/versions", true)
    xhr.send()
}

export default fetchVersions
