$(document).ready(function () {
    $("#submit-button").click(function () {
        $.ajax({
            url: "/short",
            type: "get",
            data: {
                url : $("#url-input").val(),
                customKey : $("#path-input").val(),
                customExpire : $("#time-input").val()
            },
            success: function (response) {
                console.log(response);
                $("#result-value-label").html(
                    "<a href=" + response + ">" + response + "</a>"
                );
            },
            error: function (response) {
                console.log(response);
                $("#result-value-label").html(
                    "<p> Oops! Error occurred. </p>"
                );
            }
        })
    });
});