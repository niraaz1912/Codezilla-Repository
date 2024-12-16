import { BASE_API_URL } from './constants.js';

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

const transformEndpoint = (endpoint) => `${BASE_API_URL}${endpoint}`;

const logout = () => {
    const sessionId = localStorage.getItem("sessionid");
    console.log("Session ID being sent:", sessionId); // Log session ID

    fetch(transformEndpoint("/logout"), {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ sessionid: sessionId }),
    })
    .then(response => {
        if (response.ok) {
            console.log("Logout successful");
            localStorage.removeItem("sessionid");
            window.location.replace("/login.html");
        } else {
            console.error("Logout failed:", response.statusText);
        }
    })
    .catch(err => {
        console.error("Error during logout:", err);
    });
};

// Attach logout to the global `window` object for use in dashboard.html
window.logout = logout;
