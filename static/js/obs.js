$(() => {
    console.log("loaded")
    $("#create-scene-submit").on("click", function() {
        // var saveData = $.ajax({
        //     type: 'POST',
        //     url: "/twitch/update",
        //     data: JSON.stringify({
        //         title: $("#title-text").val()
        //     }),
        //     contentType: "application/json; charset=utf-8",
        //     success: function(resultData) { alert("Save Complete") }
        // });
        var name = $("#create-scene-name").val()
        if (name.length > 0) {
            console.log("Save " + name)
        }
    });
});