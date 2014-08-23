$(document).ready(function() {
    $("#addCondition").click(function(e) {
        e.preventDefault();
        addCondition();
        updateConditions();
    });

    setupUpdateConditions();
    setupDraggables();
});

function setupUpdateConditions() {
    $("#updateConditions").click(updateConditions);
}

function setupDraggables() {
    $("#sortable").sortable();
    $("#sortable").disableSelection();
}

function updateConditions() {
    var conditions = [];
    $("#sortable li").each(function(k, v) {
        var condition = {
            Type: $(v).attr("data-type"),
            Text: $(v).text(),
            Field: $(v).attr("data-field"),
            Operator: $(v).attr("data-operator"),
            Value: $(v).attr("data-value")
        }

        conditions.push(condition);
    });

    $.ajax({
        url: "/conditions",
        method: "PUT",
        data: {
            conditions: JSON.stringify(conditions)
        }
    });
}

function addCondition() {
    if (!$("#value").val()) {
        return
    }

    var field = $("#field").val();
    var operator = $("#operator").val();
    var value = $("#value").val();

    var logicTemplate = $("#sortable .logic").first().clone();
    logicTemplate.text("AND");

    var equalityTemplate = $("#sortable .equality").first().clone();
    equalityTemplate.attr("data-field", field);
    equalityTemplate.attr("data-operator", operator);
    equalityTemplate.attr("data-value", value);
    equalityTemplate.text(field + " " + operator + " " + value);

    $("#sortable .scope").last().before(logicTemplate).before(equalityTemplate);

    $("#value").val("");
}