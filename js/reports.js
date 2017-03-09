(function () {

// Constants that map directly to the backend's Importance type.
const noChange = 'No Change'
const minorUpdate = 'Minor Update'
const contentChange = 'Content Change'
const rewritten = 'Major Rewrite'
const deleted = 'Deleted Page'

const green = 'rgb(174, 255, 139)'
const yellow = 'rgb(255, 236, 131)'
const red = 'rgb(255, 134, 149)'

$('#monitorreports').children('.reportsummary').each(function () {
    let _this = $(this)
    _this.css('background-color', red)
    let change = $(_this.find('.change')[0]).text()
    if (change.includes(noChange)) {
        _this.css('background-color', green)
    } else if (change.includes(minorUpdate) || change.includes(contentChange)) {
        _this.css('background-color', yellow)
    } else {
        _this.css('background-color', red)
    }
})

})()