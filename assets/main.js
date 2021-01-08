window.onload = function() {
    var image = document.getElementById("home_camera");

    function updateImage() {
        image.src = image.src.split('?')[0] + "?" + Math.random();
    }

    setInterval(updateImage, 1000);
}


function showClock(){
    var now = new Date();
    msg = now.getHours()+":"+now.getMinutes()+":"+now.getSeconds();
    document.getElementById("time").innerHTML = msg;
}
setInterval('showClock()', 1000);