$(() => {
    setInterval(
        function(){ 
            $.ajax({
                type: 'GET',
                url: "/avatar_status",
                success: function(resultData) {
                    console.log("here")
                    $("#mouth").addClass("open")
                    $("#mouth").removeClass("close")
                },
                error: function() {
                    console.log("bad")
                    $("#mouth").removeClass("open")
                    $("#mouth").addClass("close")
                }
            });
        }, 
        500
    );
});