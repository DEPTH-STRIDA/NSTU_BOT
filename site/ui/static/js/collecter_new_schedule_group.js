// Функция для отключения клавиши Enter
function disableEnterKey(event) {
    if (event.key === "Enter") {
      event.preventDefault();
    }
  }
  
  // Функция для сбора данных из полей в массив
  function collectData() {
    var allData = [];
  
    // Обходим элементы и собираем данные
    for (var weekType of ["even", "odd"]) {
      for (var day = 0; day < 7; day++) {
        var dataArray = [];
        for (var index = 0; ; index++) {
          var inputId = weekType + "_" + getDayName(day) + "_" + index;
          var inputElement = document.getElementById(inputId);
  
          if (!inputElement) {
            break; // Прерываем цикл, если элемент не найден
          }
  
          var inputValue = inputElement.value;
          dataArray.push(inputValue);
        }
        allData.push(dataArray);
      }
    }
  
    return allData;
  }
  
  // Функция для получения названия дня по его номеру (0 - воскресенье, 1 - понедельник и т.д.)
  function getDayName(dayNumber) {
    var daysOfWeek = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday","sunday"];
    return daysOfWeek[dayNumber];
  }
  function splitArraysToJsonString(inputArray) {
  if (inputArray.length !== 14) {
    console.error("Неверный размер входного массива. Ожидалось 14 подмассивов.");
    return null;
  }

  var evenWeek = inputArray.slice(7);
  var oddWeek = inputArray.slice(0, 7);

  var result = {
    Even_week_schedule: evenWeek,
    Odd_week_schedule: oddWeek
  };

  return JSON.stringify(result);
}
function collect(){
    var allWeekData = collectData();
    console.log("Все данные:", allWeekData);
    var jsonString = splitArraysToJsonString(allWeekData);
    return jsonString
}

function editShedule() {

  var buttoninputElement = document.getElementById("GroupName");
  var groupNameContent = buttoninputElement.textContent || buttoninputElement.innerText;


    var xhr = new XMLHttpRequest();
    var url = "/collect_schedule";

    var schedule_json=collect()
    if (schedule_json=="error"){
        alert("Неизвестная ошибка.")
        return
    }
    //alert(encodeURIComponent(schedule_json))
    var params = "schedule=" + encodeURIComponent(schedule_json)+"&id="+tg.initDataUnsafe.user.id+"&groupName="+groupNameContent;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    // Обработчик события onload для выполнения действий после успешного завершения запроса
    xhr.onload = function() {
        if (xhr.status == 200) {
            alert(xhr.responseText)
        } else {
            alert(xhr.responseText)
        }
    };

    xhr.send(params);
}