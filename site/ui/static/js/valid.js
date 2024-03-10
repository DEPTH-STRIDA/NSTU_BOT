let tg = window.Telegram.WebApp
function initialization() {
    let intervalId = setInterval(function() {
        if (tg.initData && tg.initDataUnsafe && tg.initDataUnsafe.user && tg.initDataUnsafe.user.id) {
            access(tg.initData)
            clearInterval(intervalId);
            return
        } else {
            console.error('tg.initData is not available yet.');
        }
    }, 3000)
}
function access(type_) {
    var xhr = new XMLHttpRequest();
    var url = "/validate";
    var params = "token=" + type_;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.onload = function() {
        if (xhr.status == 200) {
            window.location.href = "/home?id=" + tg.initDataUnsafe.user.id;
        } else {
            console.error("AJAX request failed with status " + xhr.status);
        }
    };

    xhr.send(params);
}