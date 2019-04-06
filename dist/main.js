/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./src/index.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./src/index.js":
/*!**********************!*\
  !*** ./src/index.js ***!
  \**********************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

eval("__webpack_require__(/*! ./styles.scss */ \"./src/styles.scss\");\n\ndocument.getElementById(\"input\").onsubmit = function (el, ev) {\n    var textArea = document.getElementsByName(\"data\")[0];\n\n    // if (el.target[0].value.indexOf(\"\\t\") >= 0) {\n    // where to put this error...\n    // }\n\n    textArea.disabled = true;\n    var encodedData = encodeURIComponent(el.target[0].value);\n\n    var request = new(XMLHttpRequest);\n    request.open(\"POST\", \"/validate\");\n    request.send(\"data=\" + encodedData);\n    request.onreadystatechange = function (ev) {\n        if (ev.target.readyState === 4) {\n            textArea.disabled = false;\n            setResults(JSON.parse(this.response));\n        }\n    }\n\n    // prevent the default behavior of navigating to the action (don't load a new page)\n    return false;\n}\n\n\n// data is a map[string][]err\nfunction setResults(data) {\n    var results = document.getElementsByClassName(\"result\");\n    var tables = document.getElementsByClassName(\"errors\");\n    var tabs = document.getElementsByClassName(\"tab\");\n\n    // Set the tab contents of each validation with highlights.\n    [].forEach.call(results, function (result) {\n        var version = result.getAttribute('data-version');\n\n        // no errors for this version\n        var errorIndex = parseInt(result.getAttribute('data-error-table-index'), 10);\n        if (data[version].length == 0) {\n            // There are no errors for this version\n            tables[errorIndex].children[1].innerHTML = \"no errors!\";\n            tabs[errorIndex].classList.remove(\"error-color\");\n            result.firstElementChild.innerHTML = document.getElementsByName(\"data\")[0].value;\n        } else {\n            // handle the errors case for this version\n            var errors = \"\";\n            data[version].forEach(function (error) {\n                errors += \"<tr><td>\" + error.Error + \"</td></tr>\";\n            });\n            tables[errorIndex].children[1].innerHTML = errors;\n            tabs[errorIndex].classList.add(\"error-color\");\n            result.firstElementChild.innerHTML = keyToRegexes(data[version][0], document.getElementsByName(\"data\")[0].value);\n        }\n    });\n}\n\n// keyToRegexes runs a series of regexes over the input to markup the document when there are validation errors.\nfunction keyToRegexes(error, value) {\n    var keys = error.Key.split(\".\");\n    var v = error.Value;\n    // each key leads to a deeper key...\n    for (var i = 0; i < keys.length - 1; i++) {\n        var reg = new RegExp(\"\\(\" + \"[ -] \".repeat(i) + keys[i] + \"\\):\");\n        value = value.replace(reg, '<span class=\"red\">$1</span>:');\n    }\n\n    // the last key will be on the same line as the value.\n    var reg = new RegExp(\"\\(\" + \"[ -] \".repeat(keys.length - 1) + keys[keys.length - 1] + \":\\\\s*\\\"?\" + v + \"\\\"?\\)\");\n    // console.log(reg);\n    value = value.replace(reg, '<span class=\"red\">$1</span>');\n    return value;\n}\n\nfunction example() {\n    document.getElementsByName(\"data\")[0].value = `apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-deployment\n  labels:\n    app: nginx\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      contaisdsners:\n      - name: nginx\n        image: nginx:1.7.9\n        ports:\n        - containerPort: 80\n`\n}\n\nfunction select(el) {\n    var results = document.getElementsByClassName(\"result\");\n    [].forEach.call(document.getElementsByClassName(\"tab\"), function (tab, i) {\n        if (tab === el) {\n            results[i].classList.remove(\"hide\");\n            tab.classList.add(\"selected\");\n            return;\n        }\n        results[i].classList.add(\"hide\");\n        tab.classList.remove(\"selected\");\n    })\n}\n\n\n//# sourceURL=webpack:///./src/index.js?");

/***/ }),

/***/ "./src/styles.scss":
/*!*************************!*\
  !*** ./src/styles.scss ***!
  \*************************/
/*! no static exports found */
/***/ (function(module, exports) {

eval("throw new Error(\"Module parse failed: Unexpected character '@' (1:0)\\nYou may need an appropriate loader to handle this file type.\\n> @charset \\\"utf-8\\\";\\n| @import \\\"~bulma/bulma\\\";\\n| \");\n\n//# sourceURL=webpack:///./src/styles.scss?");

/***/ })

/******/ });