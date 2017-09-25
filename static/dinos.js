function init() {
  getPlayerDinos();
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

function addLabelVal(parent, label, val) {
  var lblName = document.createElement("span");
  var lblNameTxt = document.createTextNode(label + ": ");
  lblName.appendChild(lblNameTxt);

  var valName = document.createElement("span");
  var valNameTxt = document.createTextNode(val);
  valName.appendChild(valNameTxt);

  var brName = document.createElement("br");

  parent.appendChild(lblName);
  parent.appendChild(valName);
  parent.appendChild(brName);
}

function makeDinoImg(dino) {
  var i = document.createElement("img");
  i.setAttribute("src", '/img/' + dino.Name.toLowerCase() + '.png');
  i.setAttribute("alt", dino.Name.toLowerCase());
  return i
}

function MakeOwnedDinoHTML(dino) {
  var parent = document.createElement("span");
  parent.appendChild(makeDinoImg(dino));
  parent.appendChild(document.createElement("br"));
  addLabelVal(parent, "Name", dino.Name);
  addLabelVal(parent, "ID", dino.ID);
  addLabelVal(parent, "Power", dino.Power);
  addLabelVal(parent, "Health", dino.Health);
  addLabelVal(parent, "Found At", '' + dino.Latitude + ',' + dino.Longitude);
  return parent
}

// TODO remove duplication
function getPlayerDinos() {
  // TODO handle unauthenticated
  var dinoParent = document.getElementById("dinoParent");
  ajax("/api/dinos", function(data) {
    while (dinoParent.firstChild) {
      dinoParent.removeChild(dinoParent.firstChild);
    }

    var dinos = JSON.parse(data);
    for (i = 0; i < dinos.length; i++) {
      var dino = dinos[i];
      var dinoSpan = MakeOwnedDinoHTML(dino);
      dinoParent.appendChild(dinoSpan);
      dinoParent.appendChild(document.createElement("br"));
    }
  });
}


init();
