$(() => {
    $("#videos").on("click", ".upload", function(){
        var id = $(this).parent().attr("id")
        var playlist_id = $(this).siblings(".select-playlist").find(":selected").val()
        var saveData = $.ajax({
            type: 'POST',
            url: "/youtube_upload",
            data: JSON.stringify({
                recording_id: parseInt(id) || -1,
                playlist_id: playlist_id || ""
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) {},
            error: function(request, status, error) {
                alert(request.responseText);
            }
        });
    });
    $("#categories").on("click", ".save-category", function() {
        var related_id = $(this).siblings(".category-options").find(":selected").val()
        var category_name = $(this).siblings(".category-name").text()
        var saveData = $.ajax({
            type: 'POST',
            url: "/youtube_category",
            data: JSON.stringify({
                related_id: related_id || "",
                category_name: category_name || ""
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) {},
            error: function(request, status, error) {
                alert(request.responseText);
            }
        });
    })
});