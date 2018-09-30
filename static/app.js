document.getElementById("input").onsubmit = function (el, ev) {
    var textArea = document.getElementsByName("data")[0];

    // if (el.target[0].value.indexOf("\t") >= 0) {
    // where to put this error...
    // }

    textArea.disabled = true;
    var encodedData = encodeURIComponent(el.target[0].value);

    var request = new(XMLHttpRequest);
    request.open("POST", "/validate");
    request.send("data=" + encodedData);
    request.onreadystatechange = function (ev) {
        if (ev.target.readyState === 4) {
            textArea.disabled = false;
            setResults(JSON.parse(this.response));
        }
    }

    // prevent the default behavior of navigating to the action (don't load a new page)
    return false;
}


// data is a map[string][]err
function setResults(data) {
    var results = document.getElementsByClassName("result");
    var tables = document.getElementsByClassName("errors");
    var tabs = document.getElementsByClassName("tab");

    // Set the tab contents of each validation with highlights.
    [].forEach.call(results, function (result) {
        var version = result.getAttribute('data-version');

        // no errors for this version
        var errorIndex = parseInt(result.getAttribute('data-error-table-index'), 10);
        if (data[version].length == 0) {
            // There are no errors for this version
            tables[errorIndex].children[1].innerHTML = "no errors!";
            tabs[errorIndex].classList.remove("error-color");
            result.firstElementChild.innerHTML = document.getElementsByName("data")[0].value;
        } else {
            // handle the errors case for this version
            var errors = "";
            data[version].forEach(function (error) {
                errors += "<tr><td>" + error.Error + "</td></tr>";
            });
            tables[errorIndex].children[1].innerHTML = errors;
            tabs[errorIndex].classList.add("error-color");
            result.firstElementChild.innerHTML = keyToRegexes(data[version][0], document.getElementsByName("data")[0].value);
        }
    });
}

// keyToRegexes runs a series of regexes over the input to markup the document when there are validation errors.
function keyToRegexes(error, value) {
    var keys = error.Key.split(".");
    var v = error.Value;
    // each key leads to a deeper key...
    for (var i = 0; i < keys.length - 1; i++) {
        var reg = new RegExp("\(" + "[ -] ".repeat(i) + keys[i] + "\):");
        value = value.replace(reg, '<span class="red">$1</span>:');
    }

    // the last key will be on the same line as the value.
    var reg = new RegExp("\(" + "[ -] ".repeat(keys.length - 1) + keys[keys.length - 1] + ":\\s*\"?" + v + "\"?\)");
    // console.log(reg);
    value = value.replace(reg, '<span class="red">$1</span>');
    return value;
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
      contaisdsners:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
`
}

function select(el) {
    var results = document.getElementsByClassName("result");
    [].forEach.call(document.getElementsByClassName("tab"), function (tab, i) {
        if (tab === el) {
            results[i].classList.remove("hide");
            tab.classList.add("selected");
            return;
        }
        results[i].classList.add("hide");
        tab.classList.remove("selected");
    })
}
