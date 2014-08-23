$(document).ready(function() {
    $("#addCondition").click(function(e) {
        e.preventDefault();

        addCondition();

        equalities = [];

        $("#sortable li").each(function(key, condition) {
            if ($(condition).attr("data-type") == "equality") {
                var condition = {
                    Field: $(condition).attr("data-field"),
                    Operator: $(condition).attr("data-operator"),
                    Value: $(condition).attr("data-value")
                }

                equalities.push(condition);
            }
        });

        updateConditions();
    });

    $("#sortable").sortable({
        change: function() {
            updateConditions();
        }
    });
    $("#sortable").disableSelection();
});

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
        url: "/updateConditions",
        method: "POST",
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