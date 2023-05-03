$(() => {
    $("#change-stream").on("click", function() {
        var saveData = $.ajax({
            type: 'POST',
            url: "/twitch/update",
            data: JSON.stringify({
                title: $("#title-text").val(),
                category_id: $("#choose-game").find(":selected").val(),
                category_name: $("#choose-game").find(":selected").text()
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) { alert("Save Complete") }
        });
    });
    $("#search-categories-submit").on("click", function(e) {
        var query = $("#search-categories").val()
        var saveData = $.ajax({
            type: 'POST',
            url: "/twitch/search/categories",
            data: JSON.stringify({
                query: query
            }),
            contentType: "application/json; charset=utf-8",
            success: function(resultData) {
                $("#categories").html("")
                var resp = JSON.parse(resultData)
                var keys = Object.keys(resp)
                for(var i = 0; i < keys.length; i++) {
                    var key = keys[i]
                    var game = resp[key]
                    game.Name = key
                    var card = buildCategoryCard(game)
                    $("#categories").append(card)
                }
            }
        });
    });
    $("#categories").on("click", ".category-search-entry", function() {
        var category_id = $(this).attr("data-id")
        var category_title = $(this).find(".title").html()
        if( $("#choose-game option[value='" + category_id + "']").length == 0 ) {
            var $option = $("<option>")
            $option.attr("value", category_id)
            $option.html(category_title)
            $("#choose-game").append($option)
        }
    })
})

function buildCategoryCard(data) {
    $card = $("<div>")
    $img = $("<img>")
    $label = $("<span>")
    $btn = $("<button>")
    $img.attr("src", data.BoxArtUrl)
    $img.attr("alt", data.Name + " box art")
    $label.attr("class", "title")
    $label.text(data.Name)
    $btn.text("Select")
    $card.append($img)
    $card.append($label)
    $card.append($btn)
    $card.attr("data-id", data.ID)
    $card.attr("class", "category-search-entry")
    return $card
}