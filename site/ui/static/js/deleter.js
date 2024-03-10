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
    }else{
        
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
