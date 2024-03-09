function callGoFunction(type_) {
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
    }
  };
  
  xhr.send(params);
}

function showConfirmationPopup(type_) {
  // Функция для отображения всплывающего окна подтверждения
  var confirmation = confirm("Вы уверены, что хотите удалить группу?");
  if (confirmation) {
    callGoFunction(type_);
  }
}
var i=0
var buttons = document.querySelectorAll("#deleteButton");
buttons.forEach(function(button) {
  i+=1
  button.addEventListener("click", function() {
    const type_ = this.getAttribute('type');
    showConfirmationPopup(type_);
  });
});
  if (i==0){
var containerElement = document.getElementById("page-contain");

// Вставляем новый HTML-код внутрь контейнера
containerElement.innerHTML = '<a class="data-card" id="noGroup">' +
                              '<h2>У вас нет ни одной группы.</h2>' +
                              '<button class="button-28" role="button" name="popup-button" onclick="button_action()">Создать группу</button>' +
                              '</a>';
  }