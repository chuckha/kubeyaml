document.getElementById("input").onsubmit = function (el, ev) {
    var textArea = document.getElementsByName("data")[0];

    // if (el.target[0].value.indexOf("\t") >= 0) {
    // where to put this error...
    // }

    textArea.disabled = true;
    var encodedData = encodeURIComponent(el.target[0].value);

    // why is it called XMLHttpRequest O_o
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
    var versionNum = ["1.8", "1.9", "1.10", "1.11", "1.12"];
    var versionIds = ["one-eight", "one-nine", "one-ten", "one-eleven", "one-twelve"];

    versionIds.forEach(function (version, i) {
        // This is pretty bad, right?
        var table = document.getElementById(version + "-errors").children[1];
        if (data[versionNum[i]].length == 0) {
            document.getElementById(version + "-button").innerText = versionNum[i] + "✅";
            table.innerHTML = "None!";
            document.getElementById(version).children.item(1).innerHTML = document.getElementsByName("data")[0].value;
        } else {
            document.getElementById(version + "-button").innerText = versionNum[i] + "❌";
            var errors = "";
            data[versionNum[i]].forEach(function (error) {
                errors += "<tr><td>" + error.Error + "</td></tr>";
            });
            document.getElementById(version).children.item(1).innerHTML = keyToRegexes(data[versionNum[i]][0], document.getElementsByName("data")[0].value);
            table.innerHTML = errors;
        }
    });
}

function showResult(item) {
    var resultDiv = document.getElementById(item);
    var versions = document.getElementsByClassName("result");

    for (var i = 0; i < versions.length; i++) {
        if (versions[i] === resultDiv) {
            versions[i].classList.remove("hide");
            continue;
        }
        versions[i].classList.add("hide")
    }
}


// keyToRegexes takes a key like a.b.c.d and returns 4 regexes
// /a:/, /  b:/, /    c:/, /      d:/ and runs each one
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
