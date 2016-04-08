var ctx;
var img;

function getLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(showPosition);
    } else { 
        x.innerHTML = "Geolocation is not supported by this browser.";
    }
}

function showPosition(position){
    console.log(position.coords);
    document.getElementById("latitude-input").value = position.coords.latitude.toString().substr(0,9);
    document.getElementById("longitude-input").value = position.coords.longitude.toString().substr(0,9);
}

window.onload = function(){
    ctx = document.getElementById("earthViz").getContext("2d");
    img = new Image();   // Create new img element
    img.onload = function(){
        ctx.drawImage(img,0,0);
    };
    img.src = '../img/earthViz.png';

    getLocation();
}
