// \todo aggregate vars into single 'world' var?
var mapDinos = {};
var map;
var serverTimeDiff = 0;
var popupInterval;
var dinosaurGetIntervalMs = 3000;

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

// queryLatlon does an API query on the given latlon, and calls f with the data
function queryLatLon(lat, lon, f) {
	ajax("/query/" + lat + "/" + lon, function(data) {
		f(JSON.parse(data));
	})
}

function addDinosaursToMap(dinosaurs, map) {
	console.log('addDinosaursToMap ' + dinosaurs.length);
	var dinosaursLen = dinosaurs.length;
	for (var i = 0; i < dinosaursLen; i++) {
		addDinosaurToMap(dinosaurs[i], map);
	}
}

// durationStr returns the hours, minutes, and seconds left in the milliseconds integer `duration`, in the format '4h20m3s'
function durationShortStr(duration) {
	var ds = duration / 1000;
	var h = Math.floor(ds / 60 / 60);
	var m = Math.floor(ds / 60 % 60);
	var s = Math.floor(ds % 60);
	var str = "";
	if(h > 0) {
		str += h + "h";
	}
	if(m > 0) {
		str += m + "m";
	}
	str += s + "s";
	return str;
}

function getDinoPopupStr(dino) {
	var dinoExpireTime = Date.parse(dino.Expiration);
	var dinoExpireFromNowMs = dinoExpireTime - Date.now() - serverTimeDiff;
	return "<b>" + dino.ID + " " + dino.Name + "</b>" + "<br \>" +
		"Time left: " + durationShortStr(dinoExpireFromNowMs);
}

function addDinosaurToMap(dino, map) {
  console.log('addDinosaurToMap ' + dino.ID);
	if(mapDinos.hasOwnProperty(dino.ID)) {
    console.log('addDinosaurToMap HAS ' + dino.ID);
		return
	}
  console.log('addDinosaurToMap NEW ' + dino.ID);

	console.log("adding " + dino.ID + " specie" + dino.Name + " expiration " + dino.Expiration);

	mapDinos[dino.ID] = dino; // TODO change to store a bool or something, if dino isn't needed
	var name = dino.Name;
	// \todo add image path to specie endpoint?
	var imagePath = "/images/" + name.toLowerCase() + ".png";

	// \todo add image dimensions to specie endpoint?
	// \todo store specie icons, for speed?
	var dinoIcon = L.icon({
		dinosaur: dino,
		iconUrl: imagePath,
    //				 shadowUrl: 'leaf-shadow.png',
		iconSize:     [64, 64], // size of the icon
    //				 shadowSize:   [5, 64], // size of the shadow
		iconAnchor:   [32, 32], // point of the icon which will correspond to marker's location
    //				 shadowAnchor: [4, 62],  // the same for the shadow
		popupAnchor:  [-5, -5] // point from which the popup should open relative to the iconAnchor
	});

	var dinoExpireTime = Date.parse(dino.Expiration);
	var dinoExpireFromNowMs = dinoExpireTime - Date.now() - serverTimeDiff;

  console.log('addDinosaurToMap dinoExpireFromNowMs ' + dinoExpireFromNowMs);
  console.log('addDinosaurToMap dinoExpireTime ' + dinoExpireTime);
  console.log('addDinosaurToMap now.getTime() ' + serverTimeDiff);
  console.log('addDinosaurToMap serverTimeDiff ' + serverTimeDiff);
	if(dinoExpireFromNowMs > 0) {
		var marker = L.marker([dino.Latitude, dino.Longitude], {icon: dinoIcon});
		var popupStr = "<b>" + dino.ID + " " + dino.Name + "</b>" + "<br \>" +
		    "Time left: " + durationShortStr(dinoExpireFromNowMs);
		marker.bindPopup(popupStr).openPopup();
		marker.getPopup().Dinosaur = dino;
		marker.addTo(map);

		window.setTimeout(function() {
		  map.removeLayer(marker);
		  delete mapDinos[dino.ID];
		  console.log("removing " + dino.ID);
		}, dinoExpireFromNowMs);
	}
}

function getDinosaursHere() {
	var latLng = map.getCenter();
	var lat = latLng.lat;
	var lon = latLng.lng;
	queryLatLon(lat, lon, function(dinosaurs) {
		addDinosaursToMap(dinosaurs, map);
	})
}

function initmap() {
	map = new L.Map('mapid');

	// create the tile layer with correct attribution
	var osmUrl='http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
	var osmAttrib='Map data © <a href="http://openstreetmap.org">OpenStreetMap</a> contributors';
	var osm = new L.TileLayer(osmUrl, {attribution: osmAttrib});

	var defaultLat = 39.7392;
	var defaultLon = -104.9903;
	var defaultZoom = 15

	map.setView(new L.LatLng(defaultLat, defaultLon), defaultZoom);
	map.addLayer(osm); 
  
	map.locate({setView: true, maxZoom: 18});


  map.on('move', getDinosaursHere);
	map.on('moveend', function() {
		console.log("map moved - calling getDinosaursHere()");
		getDinosaursHere();
	});
	map.on('popupopen', function(popupEvent) {
		var popup = popupEvent.popup;
		var dino = popup.Dinosaur;
		if(dino) {
			popup.setContent(getDinoPopupStr(dino));
			console.log('setting interval ' + dino.ID)
			popup.UpdateInterval = window.setInterval(function() {
				popup.setContent(getDinoPopupStr(dino));
			}, 1000);
		}
	})
	map.on('popupclose', function(popupEvent) {
		var popup = popupEvent.popup;
		if(popup.UpdateInterval) {
			clearInterval(popup.UpdateInterval)
			console.log('clearing interval ' + popup.Dinosaur.ID)
		}
	})
  getDinosaursHere();

}

function setLocation() {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(setPosition);
	}
}
function setPosition(position) {
	console.log('' + position.coords.latitude + ' ' + position.coords.longitude);
	map.setView(new L.LatLng(position.coords.latitude, position.coords.longitude), 18);
}

function getServerTimeDiff() {
	var startTime = Date.now();
	ajax("/now", function(timeStr) {
		var endTime = Date.now();
		var latency = endTime.getTime() - startTime.getTime();
		var serverTime = Date.parse(timeStr);
		// we subtract (latency / 2), because that's presumably how long it took to get back to us after the server created the timestamp
		serverTimeDiff = Date.now() - serverTime -  (latency / 2);
	})
}

function init() {
  getServerTimeDiff();
	initmap();
	setLocation();
	window.setInterval(getDinosaursHere, dinosaurGetIntervalMs);
}

init();
