// var serverHost = "http://<你的实际ip地址>:8080";
var serverHost = "http://127.0.0.1:8080";

function queryParams() {
    var username = localStorage.getItem("username");
    var token = localStorage.getItem("token");
    return "username=" + username + "&token=" + token;
}

function logout() {
    localStorage.removeItem("username");
    localStorage.removeItem("token");
    window.location = "/static/view/signin.html";
}

String.prototype.format = function(args) {
    var result = this;
    if (arguments.length > 0) {
        if (arguments.length == 1 && typeof args == "object") {
            for (var key in args) {
                if (args[key] != undefined) {
                    var reg = new RegExp("({" + key + "})", "g");
                    result = result.replace(reg, args[key]);
                }
            }
        } else {
            for (var i = 0; i < arguments.length; i++) {
                if (arguments[i] != undefined) {
                    var reg = new RegExp("({)" + i + "(})", "g");
                    result = result.replace(reg, arguments[i]);
                }
            }
        }
    }
    return result;
};