const isUserLoggedIn = () => {
    const sessionId = localStorage.getItem("sessionid");

    if (sessionId) {
        window.location.replace("/dashboard.html")
    }
}

const publicRoutes = ["/login.html", "/signup.html"];
const privateRoutes = ["/dashboard.html"];

window.onload = () => {
    const sessionId = localStorage.getItem("sessionid");
    if (sessionId) {
        if (publicRoutes.includes(window.location.pathname)) {
            window.location.replace("/dashboard.html")
        }
    }

    if (!sessionId) {
        if (privateRoutes.includes(window.location.pathname)) {
            window.location.replace("/login.html");
        }
    }
}