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
});