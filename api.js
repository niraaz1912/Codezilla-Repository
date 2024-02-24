const transformHeaders = (endpoint) => {
  const sessionId = localStorage.getItem("session");

  if(!sessionId) {
    window.location.replace("/login.html")
    return;
  }

  return {
    "Accept": "application.json",
    "Content-Type": "application/json",
  }
}

const isUserLoggedIn = () => {
  const sessionId = localStorage.getItem("sessionid");

  if(sessionId) {
    window.location.replace("/dashboard.html")
  }
}

const loginBtn = document.getElementById("submit-btn");

loginBtn.addEventListener("click", () => {
  const emailDiv = document.getElementById("email");
  const passwordDiv = document.getElementById("password");

  const userEmail = emailDiv.value;
  const userPwd = passwordDiv.value;

  testLogin(userEmail, userPwd);
});


function testLogin(userEmail, userPwd) {
  const requestBody = {
    username: userEmail,
    password: userPwd,
  };

  const response = fetch("http://heron.cs.umanitoba.ca:8081/" + "login", {
    method: "POST",
    headers: {
      Accept: "application.json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestBody),
  })
    .then((res) => res.json())
    .then((out) => {
      if("sessionid" in out) {
        localStorage.setItem("sessionid", out["sessionid"])
      }
    });
}
