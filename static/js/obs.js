$(() => {
    $("#create-scene-submit").on("click", function() {
        $.ajax({
            type: 'POST',
            url: "/obs/scene/create",
            data: JSON.stringify({
                name: $("#create-scene-name").val()
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) {
                var $options = $("#select-scene-options")
                $options.html("")
                for(var i = 0; i < resultData.names.length; i++) {
                    $option = $("<option>")
                    $option.val(resultData.names[i])
                    $option.text(resultData.names[i])
                    $options.append($option)
                }
            }
        });
    });
    function hexToRgb(hex) {
        var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
        return result ? {
            r: parseInt(result[1], 16),
            g: parseInt(result[2], 16),
            b: parseInt(result[3], 16),
            a: 255
        } : null;
    }
    $("#task-background-enabled").on("change", function() {
        $("#background-color-label").toggle();
        $("#background-color").toggle();
    });
    $("#record-enabled").on("change", function() {
        $("#record-output").toggle();
    });
    $("#task-submit").on("click", function() {
        var payload = {
            text: $("#task-text").val() || "",
            width: parseInt($("#task-width").val()) || 0,
            height: parseInt($("#task-height").val()) || 0,
            pos_x: parseInt($("#task-posx").val()) || 0,
            pos_y: parseInt($("#task-posy").val()) || 0,
            color: hexToRgb($("#task-color").val()) || {"r": 0, "g": 0, "b": 0, "a": 255}
        }
        if ( $("#task-background-enabled").is(":checked") ) {
            payload.background = {
                color: hexToRgb($("#background-color").val()) || {"r": 0, "g": 0, "b": 0, "a": 255}
            }
        }
        $.ajax({
            type: 'POST',
            url: "/obs/task",
            data: JSON.stringify(payload)
        });
    });
    $("#update-stream").on("click", function() {
        var recordOriginal = Boolean($("#stream-metadata").attr("data-recording")) || false
        var recordEnabled = $("#record-enabled").is(":checked")
        var payload = {
            stream: $("#stream-enabled").is(":checked"),
            record: recordEnabled
        }
        if( recordOriginal && !recordEnabled ) {
            payload.output_file = $("#record-output-name").val() || "default.mp4"
        }
        console.log(payload)
    })
});