import { BASE_API_URL } from './constants.js';

const publicRoutes = ["/login.html", "/signup.html"];
const privateRoutes = ["/dashboard.html", "user_dashboard.html"];

window.onload = () => {
    const sessionId = localStorage.getItem("sessionid");
    if (sessionId) {
        if (publicRoutes.includes(window.location.pathname)) {
            window.location.replace("/user_dashboard.html")//****************************** */
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

export const getUserRole = () => {
    const sessionId = localStorage.getItem("sessionid");

    return fetch(`${BASE_API_URL}/user/role`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "SessionID": sessionId, // Pass session ID
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Failed to fetch user role");
        }
        return response.json();
    })
    .then(data => data.role)
    .catch(err => {
        console.error("Error fetching user role:", err);
        return null;
    });
};
