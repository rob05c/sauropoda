function login() {
  console.log('login calling');
  var name = document.getElementById("name").value;
  var pass = document.getElementById("pass").value;
  console.log('login calling ' + name + ":" + pass);

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        window.location.href = '/';
        document.getElementById("loginFail").textContent = "Success! Redirecting...";
      } else {
        document.getElementById("loginFail").textContent = "Bad username or password";
      }
    }
  };
  xhttp.open("POST", '/api/login', true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send(JSON.stringify({user:name, pass:pass}));
}

function create() {
  console.log('create calling');
  var name = document.getElementById("name").value;
  var pass = document.getElementById("pass").value;
  console.log('create calling ' + name + ":" + pass);

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        window.location.href = '/';
        document.getElementById("loginFail").textContent = "Success! Redirecting...";
      } else {
        document.getElementById("loginFail").textContent = "Someone already has that name.";
      }
    }
  };
  xhttp.open("POST", '/api/createuser', true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send(JSON.stringify({user:name, pass:pass}));
}
