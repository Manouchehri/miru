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

// Set the background color of each report based on the severity
// of the change detected on the site.
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

// Attach a handler to each report that will display its contents when
// clicked a first time, and hide them the second time.
$('.moreinfo').each(function () {
    let _this = $(this)
    _this.css('visibility', 'hidden')
    _this.parent().css('height', '40px')
    let visible = false
    _this.parent().click(function (target) {
        visible = !visible
        if (visible) {
            _this.css('visibility', 'visible')
            _this.parent().css('height', 'inherit')
        } else {
            _this.css('visibility', 'hidden')
            _this.parent().css('height', '40px')
        }
    })
})

})()