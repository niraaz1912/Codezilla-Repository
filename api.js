import { BASE_API_URL } from './constants.js';


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

const transformEndpoint = (endpoint) => `${BASE_API_URL}${endpoint}`;


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


function login(username, password) {
  fetch(`${BASE_API_URL}/login`, {
      method: "POST",
      headers: {
          "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
  })
  .then(response => {
      if (response.ok) {
          return response.json();
      } else {
          console.error("Login failed with status:", response.status);
          throw new Error("Login failed");
      }
  })
  .then(data => {
      console.log("Session ID received:", data.sessionid); // Log the session ID
      localStorage.setItem("sessionid", data.sessionid);  // Store session ID in localStorage
      window.location.replace("/dashboard.html");
  })
  .catch(error => {
      console.error("Error during login:", error);
  });
}


// #endregion