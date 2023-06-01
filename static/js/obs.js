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

    $("#overlay-submit").on("click", function() {
        var payload = {
            text: $("#overlay-text").val() || "",
            text_width: parseInt($("#overlay-text-width").val()) || 0,
            text_height: parseInt($("#overlay-text-height").val()) || 0,
            text_posx: parseInt($("#overlay-text-posx").val()) || 0,
            text_posy: parseInt($("#overlay-text-posy").val()) || 0,
            text_color: hexToRgb($("#overlay-text-color").val()) || {"r": 0, "g": 0, "b": 0, "a": 255},
            background_color: hexToRgb($("#overlay-background-color").val()) || {"r": 0, "g": 0, "b": 0, "a": 255},
            enabled: $("#overlay-enabled").is(":checked")
        }
        $.ajax({
            type: 'POST',
            url: "/obs/overlay",
            data: JSON.stringify(payload)
        });
    });
    $("#avatar-text").keypress(function(event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            console.log("here")
            $("#avatar-text-submit").click();
        }
    });
    $("#avatar-text-submit").on("click", function() {
        payload = {
            text: $("#avatar-text").val()
        }
        $.ajax({
            type: 'POST',
            url: "/avatar_status",
            data: JSON.stringify(payload)
        })
        $("#avatar-text").val("")
    })
    $("#update-stream").on("click", function() {
        var streamEnabled = $("#stream-enabled").is(":checked")
        var recordEnabled = $("#record-enabled").is(":checked")
        var payload = {
            stream: streamEnabled,
            record: recordEnabled
        }
        $.ajax({
            type: 'POST',
            url: "/obs/stream",
            data: JSON.stringify(payload),
            success: function() {
                if (!streamEnabled) {
                    $("#streaming-status").removeClass("green");
                    $("#streaming-status").addClass("red")
                } else {
                    $("#streaming-status").removeClass("red");
                    $("#streaming-status").addClass("green")
                }
                if (!recordEnabled) {
                    $("#recording-status").removeClass("green");
                    $("#recording-status").addClass("red")
                } else {
                    $("#recording-status").removeClass("red");
                    $("#recording-status").addClass("green")
                }
            }
        });
    })
});