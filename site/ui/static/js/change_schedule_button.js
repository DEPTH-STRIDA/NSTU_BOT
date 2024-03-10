

function test(name) {
  //alert("Функция запущена " + name);

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