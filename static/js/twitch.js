$(() => {
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
    });
    $("#search-categories").on("input", function(e) {
        var query = $(this).val()
        var saveData = $.ajax({
            type: 'POST',
            url: "/twitch/search/categories",
            data: JSON.stringify({
                query: query
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) { console.log(resultData) }
        });
    });
})