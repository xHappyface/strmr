$(() => {
    console.log("loaded")
    $("#change-title").on("click", function() {
        var saveData = $.ajax({
            type: 'POST',
            url: "/twitch/update",
            data: JSON.stringify({
                title: $("#title-text").val()
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) { alert("Save Complete") }
        });
        console.log($("#title-text").val())
    })
})