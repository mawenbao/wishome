// load captcha
$(document).ready(function(){
    loadCaptcha();
})

function loadCaptcha() {
    $.ajax({
        type: 'POST',
        url: '/captcha/getcaptcha',
        data: {captchaid: $('#captchaid').val()},
    }).done(function(resp) {
        $('#captchaid').val(resp.id);
        $('#captchaimage').attr('src', resp.imageurl + "&v=" + new Date().getTime());
    });
}

