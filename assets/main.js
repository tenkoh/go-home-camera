function doCalibrate(){
    var q = {
        "exposure": $('[name=exposure] option:selected').val()
    }
    $.get("/api/calibration/", q).done(function(data){
        console.log(data);
    })
    confirm('調整中。約5秒後に自動的に設定を更新します。');
}

function updateImage() {
    var image = document.getElementById("home_camera");
    image.src = image.src.split('?')[0] + "?" + Math.random();
}
setInterval(updateImage, 1000);

function showClock(){
    var now = new Date();
    msg = now.getHours()+":"+now.getMinutes()+":"+now.getSeconds();
    document.getElementById("time").innerHTML = msg;
}
setInterval(showClock, 1000);

window.onload = function() {
    updateImage();
    // showClock();
}