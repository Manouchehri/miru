(function (){

console.log('Registering handlers')
const submitButtons = $('a.submitbtn')
submitButtons.click(function () {
    const _this = $(this)
    console.log('Got click for', _this)
    _this.parent().submit()
})

})()