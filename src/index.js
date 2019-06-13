require('./styles.scss');

// Makes the error tabs work
var tabs = document.getElementById("error-tabs").children;
var contents = document.getElementById('error-tab-contents').children;
for (var i=0; i<tabs.length; i++) {
    tabs.item(i).addEventListener("click", function(event) {
        // Sets the tab to be active
        for (var j=0; j<tabs.length; j++) {
            tabs.item(j).classList.remove('is-active');
        }
        event.currentTarget.classList.add('is-active');

        // Set content to display
        var dataId = event.currentTarget.dataset['tab'];
        var contentEl = document.querySelector('div[data-content="' + dataId + '"]');
        for (var j=0; j<contents.length; j++) {
            contents.item(j).classList.add('is-display-none');
        }
        contentEl.classList.remove('is-display-none');
    })
}


// keyToRegexes runs a series of regexes over the input to markup the document when there are validation errors.
function keyToRegexes(error, value) {
    console.log(error, value);
    var keys = error.Key.split(".");
    var v = error.Value;
    // each key leads to a deeper key...
    for (var i = 0; i < keys.length - 1; i++) {
        var reg = new RegExp("\(" + "[ -] ".repeat(i) + keys[i] + "\):");
        value = value.replace(reg, '<span class="has-text-danger">$1</span>:');
    }

    // the last key will be on the same line as the value.
    var reg = new RegExp("\(" + "[ -] ".repeat(keys.length - 1) + keys[keys.length - 1] + ":\\s*\"?" + v + "\"?\)");
    // console.log(reg);
    value = value.replace(reg, '<span class="has-text-danger">$1</span>');
    return value;
}

function setResults(data) {
    console.log(data);
    for (var version in data) {
        var tabEl = document.querySelector('li[data-tab="'+ version +'"]');
        var contentEl = document.querySelector('div[data-content="' + version + '"]');

        if (data[version].length === 0) {
            contentEl.firstElementChild.innerHTML = "<p>✅ No errors</p>";
            contentEl.lastElementChild.innerHTML = "";
            contentEl.lastElementChild.classList.add('is-invisible');
            tabEl.classList.add('has-background-success');
            tabEl.classList.remove('has-background-danger');
        } else {
            var errors = "<ul>";
            // TODO: sort the errors, probably on the backend.
            data[version].forEach(function (error) {
                errors += "<li>" + error.Error + "</li>";
            });
            contentEl.firstElementChild.innerHTML = errors;
            tabEl.classList.add('has-background-danger');
            tabEl.classList.remove('has-background-success');
            contentEl.lastElementChild.classList.remove('is-invisible');
            contentEl.lastElementChild.innerHTML = keyToRegexes(data[version][0], document.getElementsByName("data")[0].value);
        }
    }
}


function example() {
    document.getElementsByName("data")[0].value = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
`
}
document.getElementById("example").addEventListener("click", example);
document.getElementById('input').addEventListener('submit', function(evt) {
    evt.preventDefault();
    var textArea = document.getElementsByName("data")[0];
    if (textArea.value === "") {
        // TODO blink the text area?
        return
    }

    // if (el.target[0].value.indexOf("\t") >= 0) {
    // where to put this error...
    // }

    textArea.disabled = true;
    var encodedData = encodeURIComponent(textArea.value);

    var request = new(XMLHttpRequest);
    request.open("POST", getBackendUrl() + "/validate");
    request.send("data=" + encodedData);
    request.onreadystatechange = function (ev) {
        if (ev.target.readyState === 4) {
            textArea.disabled = false;
            setResults(JSON.parse(this.response));
        }
    }

    // prevent the default behavior of navigating to the action (don't load a new page)
    return false;
});

function getBackendUrl() {
    if (window.location.protocol === 'file:') {
        // dev version assumes CORS is enabled
        return 'http://localhost:9000';
    }
    return '';
}