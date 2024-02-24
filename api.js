const loginBtn = document.getElementById("submit-btn");

loginBtn.addEventListener("click", () => {
  const emailDiv = document.getElementById("email");
  const passwordDiv = document.getElementById("password");

  const userEmail = emailDiv.value;
  const userPwd = passwordDiv.value;

  testLogin(userEmail, userPwd);
});

function login(name, password) {
  const xhr = new XMLHttpRequest();
  xhr.open("POST", "/login", true);

  // xhr.onload =
  // (_) => {
  //     if(xhr.readyState === 4)
  // }

  obj = { username: name, password: password };
  xhr.send(JSON.stringify(obj));
}

function testLogin(userEmail, userPwd) {
  const requestBody = {
    username: userEmail,
    password: userPwd,
  };

  console.log("REACHED HERE");

  const response = fetch("http://localhost:8081/login", {
    method: "POST",
    headers: {
      Accept: "application.json",
      "Content-Type": "application/json",
    },
    body: requestBody,
  })
    .then((res) => res.json())
    .then((out) => console.log(out));
}
