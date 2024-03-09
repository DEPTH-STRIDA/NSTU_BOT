let tg = window.Telegram.WebApp
function start() {
    //window.location.href = "/home?id=878413772";
    // Проверка, что tg.initData существует и не является undefined или null
    let intervalId = setInterval(function() {
        if (tg.initData && tg.initDataUnsafe && tg.initDataUnsafe.user && tg.initDataUnsafe.user.id) {
            // Вызываем функцию start только если tg.initData уже загружена
            callGoFunction(tg.initData)
            // Очистить интервал, так как данные загружены
            clearInterval(intervalId);
            return
        } else {
            console.error('tg.initData is not available yet.');
        }
    }, 3000)
}

function callGoFunction(type_) {
    var xhr = new XMLHttpRequest();
    var url = "/validate";
    var params = "token=" + type_;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    // Обработчик события onload для выполнения действий после успешного завершения запроса
    xhr.onload = function() {
        if (xhr.status == 200) {
            window.location.href = "/home?id=" + tg.initDataUnsafe.user.id;
        } else {
            console.error("AJAX request failed with status " + xhr.status);
        }
    };

    xhr.send(params);
}