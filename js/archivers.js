(function (){

const rejectButtons = $('a.submitbtn')
rejectButtons.click(function () {
    const _this = $(this)
    _this.parent().submit()
})

})