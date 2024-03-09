$(document).ready(function() {
    // Счетчики для индексов полей ввода
    let even_mondayIndex = 1;
    let even_tuesdayIndex = 1;
    let even_wednesdayIndex = 1;
    let even_thursdayIndex = 1;
    let even_fridayIndex = 1;
    let even_saturdayIndex = 1;
    let even_sundayIndex = 1;

    let odd_mondayIndex = 1;
    let odd_tuesdayIndex = 1;
    let odd_wednesdayIndex = 1;
    let odd_thursdayIndex = 1;
    let odd_fridayIndex = 1;
    let odd_saturdayIndex = 1;
    let odd_sundayIndex = 1;
    // Функция для создания нового поля ввода
    function createInput(day, index) {
        const newRow = $("<tr>");
        newRow.append(`<td><input class="input_class" type="text" onkeydown="disableEnterKey(event)" placeholder="Предмет" id="${day}_${index}" /></td>`);
        $("#"+day+"_plus_button").closest("tr").before(newRow);
    }

    // Обработчики событий для понедельника
    $("#even_monday_plus_button").on("click", function() {
        createInput("even_monday", even_mondayIndex);
        even_mondayIndex++;
    });
    $("#even_monday_minus_button").on("click", function() {
        if (even_mondayIndex > 1) {
            even_mondayIndex--;
            $("#even_monday_" + even_mondayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для вторника
    $("#even_tuesday_plus_button").on("click", function() {
        createInput("even_tuesday", even_tuesdayIndex);
        even_tuesdayIndex++;
    });
    $("#even_tuesday_minus_button").on("click", function() {
        if (even_tuesdayIndex > 1) {
            even_tuesdayIndex--;
            $("#even_tuesday_" + even_tuesdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для среды
    $("#even_wednesday_plus_button").on("click", function() {
        createInput("even_wednesday", even_wednesdayIndex);
        even_wednesdayIndex++;
    });
    $("#even_wednesday_minus_button").on("click", function() {
        if (even_wednesdayIndex > 1) {
            even_wednesdayIndex--;
            $("#even_wednesday_" + even_wednesdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для четверга
    $("#even_thursday_plus_button").on("click", function() {
        createInput("even_thursday", even_thursdayIndex);
        even_thursdayIndex++;
    });
    $("#even_thursday_minus_button").on("click", function() {
        if (even_thursdayIndex > 1) {
            even_thursdayIndex--;
            $("#even_thursday_" + even_thursdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для пятницы
    $("#even_friday_plus_button").on("click", function() {
        createInput("even_friday", even_fridayIndex);
        even_fridayIndex++;
    });
    $("#even_friday_minus_button").on("click", function() {
        if (even_fridayIndex > 1) {
            even_fridayIndex--;
            $("#even_friday_" + even_fridayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для субботы
    $("#even_saturday_plus_button").on("click", function() {
        createInput("even_saturday", even_saturdayIndex);
        even_saturdayIndex++;
    });
    $("#even_saturday_minus_button").on("click", function() {
        if (even_saturdayIndex > 1) {
            even_saturdayIndex--;
            $("#even_saturday_" + even_saturdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для воскресенья
    $("#even_sunday_plus_button").on("click", function() {
        createInput("even_sunday", even_sundayIndex);
        even_sundayIndex++;
    });
    $("#even_sunday_minus_button").on("click", function() {
        if (even_sundayIndex > 1) {
            even_sundayIndex--;
            $("#even_sunday_" + even_sundayIndex).closest("tr").remove();
        }
    });
    ///////////////////////
    // Обработчики событий для понедельника
    $("#odd_monday_plus_button").on("click", function() {
        createInput("odd_monday", odd_mondayIndex);
        odd_mondayIndex++;
    });
    $("#odd_monday_minus_button").on("click", function() {
        if (odd_mondayIndex > 1) {
            odd_mondayIndex--;
            $("#odd_monday_" + odd_mondayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для вторника
    $("#odd_tuesday_plus_button").on("click", function() {
        createInput("odd_tuesday", odd_tuesdayIndex);
        odd_tuesdayIndex++;
    });
    $("#odd_tuesday_minus_button").on("click", function() {
        if (odd_tuesdayIndex > 1) {
            odd_tuesdayIndex--;
            $("#odd_tuesday_" + odd_tuesdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для среды
    $("#odd_wednesday_plus_button").on("click", function() {
        createInput("odd_wednesday", odd_wednesdayIndex);
        odd_wednesdayIndex++;
    });
    $("#odd_wednesday_minus_button").on("click", function() {
        if (odd_wednesdayIndex > 1) {
            odd_wednesdayIndex--;
            $("#odd_wednesday_" + odd_wednesdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для четверга
    $("#odd_thursday_plus_button").on("click", function() {
        createInput("odd_thursday", odd_thursdayIndex);
        odd_thursdayIndex++;
    });
    $("#odd_thursday_minus_button").on("click", function() {
        if (odd_thursdayIndex > 1) {
            odd_thursdayIndex--;
            $("#odd_thursday_" + odd_thursdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для пятницы
    $("#odd_friday_plus_button").on("click", function() {
        createInput("odd_friday", odd_fridayIndex);
        odd_fridayIndex++;
    });
    $("#odd_friday_minus_button").on("click", function() {
        if (odd_fridayIndex > 1) {
            odd_fridayIndex--;
            $("#odd_friday_" + odd_fridayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для субботы
    $("#odd_saturday_plus_button").on("click", function() {
        createInput("odd_saturday", odd_saturdayIndex);
        odd_saturdayIndex++;
    });
    $("#odd_saturday_minus_button").on("click", function() {
        if (odd_saturdayIndex > 1) {
            odd_saturdayIndex--;
            $("#odd_saturday_" + odd_saturdayIndex).closest("tr").remove();
        }
    });

    // Обработчики событий для воскресенья
    $("#odd_sunday_plus_button").on("click", function() {
        createInput("odd_sunday", odd_sundayIndex);
        odd_sundayIndex++;
    });
    $("#odd_sunday_minus_button").on("click", function() {
        if (odd_sundayIndex > 1) {
            odd_sundayIndex--;
            $("#odd_sunday_" + odd_sundayIndex).closest("tr").remove();
        }
    });
});
