// set form autofocus

$(document).ready(function() {
    setFocusInput();
})

function setFocusInput() {
    if (!$("#inputName").val()) {
        $("#inputName").focus();
    } else {
        if (0 == $("#inputEmail").length) {
            // signin page
            $("#inputPasswd").focus();
        } else {
            if (!$("#inputEmail").val()) {
                // reset pass & signup pages
                $("#inputEmail").focus();
            } else {
                if (0 == $("#inputPasswd").length) {
                    // reset pass page
                    $("#inputCaptcha").focus();
                } else {
                    // signup page
                    $("#inputPasswd").focus();
                }
            }
        }
    }
}

