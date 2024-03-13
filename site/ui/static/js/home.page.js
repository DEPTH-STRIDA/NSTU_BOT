//////////////////////////////////////////////
///              home_script               ///
//////////////////////////////////////////////
let tg = window.Telegram.WebApp
function start() {

   // Вставляем новый HTML-код внутрь контейнер
   var containerElement = document.getElementById("page-contain");
   containerElement.innerHTML = '<a class="data-card" id="noGroup">' +
       '<h2>У вас нет ни одной группы.</h2>' +
       '<button class="button-28" role="button" name="popup-button" onclick="button_action()">Создать группу</button>' +
       '</a>';

    let intervalId = setInterval(function() {
        if (tg.initData && tg.initDataUnsafe && tg.initDataUnsafe.user && tg.initDataUnsafe.user.id) {
            var i = 0
            var buttons = document.querySelectorAll("#deleteButton");
            buttons.forEach(function(button) {
                i += 1
                button.addEventListener("click", function() {
                    const type_ = this.getAttribute('type');
                    showConfirmationPopup(type_);
                });
            });
            if (i === 0) {
                // Вставляем новый HTML-код внутрь контейнер
                var containerElement = document.getElementById("page-contain");
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
            Telegram.WebApp.onEvent('backButtonClicked', function() {
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
function chekData() {
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
function button_action(){
  window.location.href = "/new-group"
}

function button_tg() {
    tg.MainButton.setText("Создать группу");
    tg.MainButton.textColor = "#FFFFFF";
    tg.MainButton.color = "#753BBD";

    Telegram.WebApp.onEvent('mainButtonClicked', function() {
        window.location.href = "/new-group"
    });

    tg.MainButton.show()
}
//////////////////////////////////////////////
///                deleter                 ///
//////////////////////////////////////////////
function initialization(type_) {
    // Функция для вызова удаления группы
    var xhr = new XMLHttpRequest();
    var url = "/deleteGroup";
    var params = "groupname=" + type_;
    
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    
    xhr.onload = function() {
      if (xhr.readyState == 4 && xhr.status == 200) {
        var response = xhr.responseText;
        console.log(response);
        // Удаляем элемент с id "type_"
        var element = document.getElementById("id_" + type_);
        if (element) {
          element.remove();
        }
      }else{
          
      }
    };
    
    xhr.send(params);
  }
  
  function showConfirmationPopup(type_) {
    // Функция для отображения всплывающего окна подтверждения
    var confirmation = confirm("Вы уверены, что хотите удалить группу?");
    if (confirmation) {
      initialization(type_);
    }
  }
  
//////////////////////////////////////////////
///        change_schedule_button          ///
//////////////////////////////////////////////


  function test(name) {
   // Создание формы
    var form = document.createElement('form');
    form.method = 'post';
    form.action = '/changeSchedule';
  
    // Добавление поля id в форму
    var idField = document.createElement('input');
    idField.type = 'hidden';
    idField.name = 'id';
    idField.value = tg.initDataUnsafe.user.id;
  
    // Добавление других полей, если необходимо
    var scheduleField = document.createElement('input');
    scheduleField.type = 'hidden';
    scheduleField.name = 'groupName';
    scheduleField.value = (name);
  
    form.appendChild(idField);
    form.appendChild(scheduleField);
  
    // Добавление формы в body и ее автоматическая отправка
    document.body.appendChild(form);
  
    // Подписка на событие onload для обработки ответа после отправки формы
    form.onload = function() {
      alert("Функция завершила работу " + name);
    };
  
    form.submit();
  }
  
  function change_schedule(type_) {
    // Функция для вызова удаления группы
    var xhr = new XMLHttpRequest();
    var url = "/deleteGroup";
    var params = "groupname=" + type_;
  
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  
    xhr.onload = function() {
      if (xhr.readyState == 4 && xhr.status == 200) {
        window.location.href = "/home?id=" + tg.initDataUnsafe.user.id;
      }
    };
  
    xhr.send(params);
  }