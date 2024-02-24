const transformHeaders = (endpoint) => {
  const sessionId = localStorage.getItem("sessionid");

  if (!sessionId) {
    window.location.replace("/login.html")
    return;
  }

  return {
    "Accept": "application.json",
    "Content-Type": "application/json",
  }
}

const transformEndpoint = (endpoint) => {
  return "http://heron.cs.umanitoba.ca:8081/" + endpoint;
}


const postRequest = (endpoint, requestBody) => {
  return fetch(transformEndpoint(endpoint), {
    method: "POST",
    headers: {
      Accept: "application.json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestBody),
  }).catch(res => {
    if (res.status === 401) {
      window.location.replace("/login.html");
    }
  })
}

const getRequest = (endpoint) => {
  return fetch(transformEndpoint(endpoint)).catch(res => {
    if (res.status === 401) {
      window.localStorage.replace("/login.html");
    }
  })
}

// #region SIGN UP PORTION
const signUpForm = document.getElementById("sign-up-form");

signUpForm?.addEventListener("submit", (e) => {
  e.preventDefault();

  const emailDiv = document.getElementById("email");
  const passwordDiv = document.getElementById("password");

  signUp(emailDiv.value, passwordDiv.value);
})

const signUp = (userEmail, userPwd) => {
  const requestBody = {
    username: userEmail,
    password: userPwd,
  };

  const response = postRequest("/login/new", requestBody)
    .then((res) => res.json())
    .then((out) => {
      localStorage.setItem("sessionid", "User has signed up")
      window.location.replace("/dashboard.html")
    });
}
// #endregion


// #region LOGIN PORTION
const loginForm = document.getElementById("login-form");

loginForm?.addEventListener("submit", (e) => {
  e.preventDefault();

  const emailDiv = document.getElementById("email");
  const passwordDiv = document.getElementById("password");

  login(emailDiv.value, passwordDiv.value);
})


function login(userEmail, userPwd) {
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
      localStorage.setItem("sessionid", out["sessionid"])
      window.location.replace("/dashboard.html")
    });
}
// #endregion