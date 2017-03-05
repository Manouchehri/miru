(function () {

// Constants that map directly to the backend's Importance type.
const noChange = 'No Change'
const minorUpdate = 'Minor Update'
const contentChange = 'Content Change'
const rewritten = 'Major Rewrite'
const deleted = 'Deleted Page'

const green = '#00ff00'
const yellow = '#ffff00'
const red = '#ff0000'

// Color each report's background either:
// 1. Green if the change to the site is minor
// 2. Yellow if there was a change worth inspecting
// 3. Red if a dramatic rewrite or deletion was detected
let reports = $('#monitorreports')
for (let report of reports.children()) {
    console.log('Inspecting report', report)
    let rows = reports.children('tbody tr')
    console.log('Found rows', rows)
    let siteChangeRows = rows.filter((row) => {
        console.log('Inspecting row', row)
        return row
            .children('td.reportKey')
            .text()
            .includes('Site change')
    })
    console.log('Found site change rows', siteChangeRows)
    if (siteChangeRows.length === 0) {
        continue
    }
    let siteChange = siteChangeRows[0]
        .children('td.reportValue')
        .text()
    console.log('Found site change', siteChange)
    if (siteChange.includes(noChange) || siteChange.includes(minorUpdate)) {
        report.style('background-color', green)
    } else if (siteChange.includes(contentUpdate)) {
        report.style('background-color', yellow)
    } else {
        report.style('background-color', red)
    }
}

})()