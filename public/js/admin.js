// load timer results

$(document).ready(function() {
    // load timer results
    loadTimerResults();
})

function loadTimerResults() {
    setInterval(queryTimerResults, 1000);
}

function queryTimerResults() {
    $.ajax({
        type: 'POST',
        data: "sort=1",
        url: '/admin/gettimerresults'
    }).done(setTimerResultsTable);
}

function setTimerResultsTable(resp) {
    var tableBody = $('#timerResultsTableBody');
    tableBody.empty();

    $.each(resp, function(i, item) {
        var $tr = $('<tr>').append(
            $('<td>').text(item.action),
            $('<td>').text(item.avgtime),
            $('<td>').text(item.hit)
        ).appendTo(tableBody);
    });
}

