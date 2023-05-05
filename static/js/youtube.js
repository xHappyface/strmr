$(() => {
    console.log("loaded")
    $("#videos").on("click", ".upload", function(){
        var id = $(this).parent().attr("id")
        var saveData = $.ajax({
            type: 'POST',
            url: "/youtube_upload",
            data: JSON.stringify({
                recording_id: parseInt(id) || -1
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) {},
            error: function(request, status, error) {
                alert(request.responseText);
            }
        });
    })
});