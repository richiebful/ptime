var ctx;
var img;
var times = {fajr: "{{index timeMap fajr}}", dhuhr: "{{index timeMap dhuhr}}", asr: "{{.asr}}", maghrib: "{{.maghrib}}", isha: "{{.isha}}", updateFlag: true};

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

function drawTimes(){
    ctx.fillStyle = "white";
    ctx.font = "20px serif";
    ctx.fillText(times.fajr, 50, 535);
    ctx.fillText(times.dhuhr, 60, 210);
    ctx.fillText(times.asr, 550, 70);
    ctx.fillText(times.maghrib, 880, 210);
    ctx.fillText(times.isha, 900, 535);
}

window.onload = function(){
    ctx = document.getElementById("earthViz").getContext("2d");
    img = new Image();   // Create new img element
    img.onload = function(){
        ctx.drawImage(img,0,0);
	drawTimes();
    };
    img.src = '/www/img/earthViz.png';
    var d = new Date();
    console.log(-1 * d.getTimezoneOffset()/ 60.0);
    document.getElementById("date-input").value = d.getFullYear() + "/" + (d.getMonth() + 1) + "/" + d.getDate();
    document.getElementById("zone-input").value = (-1 * d.getTimezoneOffset()/60.0);
    getLocation();
}
