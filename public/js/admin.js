// load timer results

$(document).ready(function() {
    loadTimerResults();
})

function loadTimerResults() {
    setInterval(
        function(){
            $.ajax({
                type: 'GET',
                url: '/admin/gettimerresults'
            }).done(setTimerResultsTable)
        },
        1000
    )
}

function setTimerResultsTable(resp) {
    var tableBody = $('#timerResultsTBody');
    tableBody.empty();

    $.each(resp, function(i, item) {
        var $tr = $('<tr>').append(
            $('<td>').text(item.action),
            $('<td>').text(item.avgtime),
            $('<td>').text(item.hit)
        ).appendTo(tableBody);
    });
}

