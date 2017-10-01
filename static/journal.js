function init() {
  getPlayerSpecies();
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

function MakeSpeciesHTMLStr(s) {
  if(s.Name.startsWith('?')) {
    return MakeUnknownSpeciesHTMLStr(s);
  }
  return '<span>' +
    '  <img src="/img/'+s.Name.toLowerCase()+'.png" alt="'+s.Name.toLowerCase()+'"><br>' +
    '  <span>'+s.Name+'</span><br>' +
    '  <span>Height:</span><span>'+s.HeightMetres+' metres</span><br>' +
    '  <span>Length:</span><span>'+s.LengthMetres+' metres</span><br>' +
    '  <span>Weight:</span><span>'+s.WeightKg+' kilograms</span><br>' +
    '</span>';
}

function MakeUnknownSpeciesHTMLStr(s) {
  return '<span>' +
    '  <br><span>'+s.Name+'</span><br><br>' +
    '</span>';
}


// TODO remove duplication
function getPlayerSpecies() {
  // TODO handle unauthenticated
  var dinoParent = document.getElementById("dinoParent");
  ajax("/api/journal", function(data) {
    console.log("api/journal data: " + data);
    while (dinoParent.firstChild) {
      dinoParent.removeChild(dinoParent.firstChild);
    }
    var species = JSON.parse(data);
    var speciesHTMLStr = '';
    for (i = 0; i < species.length; i++) {
      var specie = species[i];
      speciesHTMLStr += MakeSpeciesHTMLStr(specie) + '<br>';
    }
    dinoParent.innerHTML = speciesHTMLStr;
  });
}


init();
