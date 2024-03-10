//document.getElementById("body").style.display = "none";
let tg = window.Telegram.WebApp
function start() {
  let intervalId = setInterval(function() {
      if (tg.initData && tg.initDataUnsafe && tg.initDataUnsafe.user && tg.initDataUnsafe.user.id) {
          var i=0
var buttons = document.querySelectorAll("#deleteButton");
buttons.forEach(function(button) {
  i+=1
  button.addEventListener("click", function() {
    const type_ = this.getAttribute('type');
    showConfirmationPopup(type_);
  });
});
  if (i===0){
var containerElement = document.getElementById("page-contain");

// Вставляем новый HTML-код внутрь контейнера
containerElement.innerHTML = '<a class="data-card" id="noGroup">' +
                              '<h2>У вас нет ни одной группы.</h2>' +
                              '<button class="button-28" role="button" name="popup-button" onclick="button_action()">Создать группу</button>' +
                              '</a>';
  }
          var scheduleButtons = document.querySelectorAll("#scheduleButton");

scheduleButtons.forEach(function(button) {
  button.addEventListener("click", function() {
    const type_ = this.getAttribute('type');
    test(type_);
  });
});
Telegram.WebApp.onEvent('backButtonClicked', function(){
        tg.close()
    });
    tg.BackButton.show()

          chekData()
          clearInterval(intervalId);
          return
      } else {
          console.error('tg.initData is not available yet.');
      }
  }, 100)
}
// function callGoFunctionData(type_) {
//   var xhr = new XMLHttpRequest();
//    var url = "/validate";
//    var params = "token=" + type_;
//    xhr.open("POST", url, true);
//    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

//  // Обработчик события onload для выполнения действий после успешного завершения запроса
//   xhr.onload = function() {
//        if (xhr.status == 200) {
//            var response = xhr.responseText;
//            console.log(response);
//        } else { 
//          var bodyElement = document.body;
//          bodyElement.remove();
//            console.error("AJAX request failed with status " + xhr.status);
//        }
//    };

//    xhr.send(params);
// }
function chekData(){
  var xhr = new XMLHttpRequest();
  var url = "/validate";
  var params = "token=" + tg.initData;
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

  // Обработчик события onload для выполнения действий после успешного завершения запроса
  xhr.onload = function() {
      if (xhr.status == 200) {
        button_tg()
        document.getElementById("body").style.display = "block";
      } else {
        console.error("AJAX request failed with status " + xhr.status);
      }
  };

  xhr.send(params);
}

function button_tg(){
  tg.MainButton.setText("Создать группу"); 
  tg.MainButton.textColor = "#FFFFFF"; 
  tg.MainButton.color = "#753BBD"; 
  
  Telegram.WebApp.onEvent('mainButtonClicked', function(){
    window.location.href = "/new-group"
  });

  tg.MainButton.show()
}
