function init() {
  getPlayerInfo();
}

function ajax(path, f) {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
		  f(this.responseText)
		}
	};
	xhttp.open("GET", path, true);
	xhttp.send();
}

function getPlayerInfo() {
  // TODO handle unauthenticated
  ajax("/api/player", function(data) {
    var playerInfo = JSON.parse(data);
    document.getElementById("player").textContent = playerInfo.name;
		withLocation(function(position) {
			document.getElementById("location").textContent = position.coords.latitude + ", " + position.coords.longitude;
		})
  });
}

function withLocation(f) {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(f);
	}
}


init();
