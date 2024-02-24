const isUserLoggedIn = () => {
    const sessionId = localStorage.getItem("sessionid");

    if (sessionId) {
        window.location.replace("/dashboard.html")
    }
}

const publicRoutes = ["/login.html", "/signup.html"];

window.onload = () => {
    const sessionId = localStorage.getItem("sessionid");
    if (sessionId && publicRoutes.includes(window.location.pathname)) {
        window.location.replace("/dashboard.html")
    }

    if (window.location.pathname !== "/login.html") {

        if (sessionId) {
            window.location.replace("/dashboard.html");
        } else {
            window.location.replace("/login.html")
        }
    }
}