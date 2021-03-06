// load timer results

$(document).ready(function() {
    queryTimerResults();
    // load timer results
    loadTimerResults();
})

function toggleSortField(elem) {
    $('.timer-sort-icon').each(function(i, item) {
        if (elem == item) {
            // set sort field
            $('#timer-sort-field').val(i);
            // set sort order
            var orderObj = $('#timer-sort-order')
            if ($(item).hasClass('fa-sort-asc')) {
                orderObj.val('0');
            } else {
                orderObj.val('1');
            }
            // set table header icon
            if (0 == orderObj.val()) {
                $(item).attr('class', 'fa fa-sort-desc timer-sort-icon');
            } else {
                $(item).attr('class', 'fa fa-sort-asc timer-sort-icon');
            }
            // query new ordered timer results
            queryTimerResults(); 
        } else {
            // reset other table header icon
            $(item).attr('class', 'fa fa-sort timer-sort-icon');
        }
    });
}

function loadTimerResults() {
    setInterval(queryTimerResults, 3000);
}

function queryTimerResults() {
    $.ajax({
        type: 'POST',
        data: {
            sortField: $('#timer-sort-field').val(),
            sortOrder: $('#timer-sort-order').val()
        },
        url: '/admin/gettimerresults'
    }).done(setTimerResultsTable);
}

function newTimerResultsRow(action, avgtime, hit) {
    return $('<tr>').append(
            $('<td>').text(action),
            $('<td>').text(avgtime),
            $('<td>').text(hit)
    );
}

function setTimerResultsTable(resp) {
    var tableBody = $('#timerResultsTableBody');
    tableBody.empty();

    // set table content
    var totalHitCount = 0; 
    var totalHitExTimer = 0;
    $.each(resp, function(i, item) {
        var currHit = parseInt(item.hit)
        totalHitCount += currHit;
        if ('admin.gettimerresults' != item.action) {
            totalHitExTimer += currHit;
        }
        newTimerResultsRow(item.action, item.avgtime, item.hit).appendTo(tableBody);
    });
    newTimerResultsRow('TOTAL (no timer)', '-', totalHitExTimer).appendTo(tableBody);
    newTimerResultsRow('TOTAL', '-', totalHitCount).appendTo(tableBody);
}

