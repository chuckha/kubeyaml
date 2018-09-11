document.getElementById("input").onsubmit = function (el, ev) {
    console.log(el, ev);
    var encodedData = encodeURIComponent(el.target[0].value);

    // why is it called XMLHttpRequest O_o
    var request = new(XMLHttpRequest);

    request.open("POST", "/validate");
    request.send("data=" + encodedData);
    request.onreadystatechange = function (ev) {
        if (ev.target.readyState === 4) {
            console.log(ev);
            console.log(this.response);
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
        var aggTableData = document.getElementById(version);

        // This is pretty bad, right?
        var table = document.getElementById(version + "-errors").children[1];
        if (data[versionNum[i]].length == 0) {
            aggTableData.innerText = "✅";
            table.innerHTML = "None!";
        } else {
            aggTableData.innerText = "❌";
            var errors = "";
            data[versionNum[i]].forEach(function (error) {
                errors += "<tr><td>" + error.Error + "</td></tr>";
            });
            table.innerHTML = errors;
        }
    });
}
