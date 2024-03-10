let tg = window.Telegram.WebApp
function callGoFunction() {
    var inputElement = document.getElementById("groupNameInput");
    var inputValue = inputElement.value;
    inputValue=inputValue.trim()
    inputElement.value=inputValue

    var xhr = new XMLHttpRequest();
    var url = "/checkGroupName";
    var params = "groupName=" + encodeURIComponent(inputValue)+"&id="+tg.initDataUnsafe.user.id;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    if (inputValue==""){
        alert("Поле пустое")
        return
    }
    // Обработчик события onload для выполнения действий после успешного завершения запроса
    xhr.onload = function() {
        if (xhr.status == 200) {
            alert(xhr.responseText)
            inputElement.disabled =true
            var buttoninputElement = document.getElementById("first_button");
            buttoninputElement.disabled =true
            change_acces()

        } else {
            alert(xhr.responseText)
        }
    };

    xhr.send(params);
}
function disableEnterKey(event) {
    if (event.key === "Enter") {
        event.preventDefault();
        callGoFunction()
    }
  }
function change_acces(){
    var block = document.getElementById("second_login");
    block.style.display = "block";
}
function button_back_tg(){
    Telegram.WebApp.onEvent('backButtonClicked', function(){
        goBack(tg.initData)
    });
    tg.BackButton.show()

}
function goBack(type_) {
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
button_back_tg()