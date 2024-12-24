import { BASE_API_URL } from './constants.js';

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
      window.location.replace("/login.html")
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
      if (!response.ok) {
          throw new Error("Login failed");
      }
      return response.json();
  })
  .then(data => {
      console.log("Session ID received:", data.sessionid);
      console.log("User role:", data.role);

      localStorage.setItem("sessionid", data.sessionid); // Save session ID
      localStorage.setItem("role", data.role); // Save user role

      // Redirect based on role
      if (data.role === "admin") {
          window.location.replace("/dashboard.html");
      } else {
          window.location.replace("/user_dashboard.html");
      }
  })
  .catch(error => {
      console.error("Error during login:", error);
      alert("Invalid login credentials.");
  });
}


let map; // Global map object
let markers = []; // Array to store markers

// Function to initialize or update the map
export function initMap(centerLat, centerLng) {
  console.log("Map Initialized")
    if (!map) {
        // Initialize the map only once
        map = new google.maps.Map(document.getElementById("map"), {
          center: { lat: 40.730610, lng: -73.935242 },
            zoom: 12,
            mapId: "map-id", 
        });

        fetchUserLocations()
        
    } else {
        // Update the map's center
        map.setCenter({ lat: centerLat, lng: centerLng });
    }
}

// Function to clear existing markers
function clearMarkers() {
    markers.forEach(marker => marker.setMap(null));
    markers = [];
}

// Function to add markers for all users
function addMarkers(userLocations) {
  clearMarkers(); // Clear old markers before adding new ones
  const bounds = new google.maps.LatLngBounds();

  for (const [email, locations] of Object.entries(userLocations)) {
      locations.forEach(location => {
          console.log(`Adding advanced marker at (${location.latitude}, ${location.longitude}) for user: ${email}`);

          const position = { lat: location.latitude, lng: location.longitude };

          // Create Advanced Marker
          const marker = new google.maps.marker.AdvancedMarkerElement({
              position: position,
              map: map, // Attach marker to the map
              title: `User: ${email}`,
          });

          markers.push(marker); // Store the marker
          bounds.extend(position); // Extend map bounds
      });
  }

  // Fit map to include all markers
  map.fitBounds(bounds);
  console.log("Map bounds updated.");
}

export async function fetchSessions(filters = {}) {
  const sessionId = localStorage.getItem("sessionid");

  if (!sessionId) {
      console.error("Session ID missing");
      return [];
  }

  // Build query string from filters
  const queryParams = new URLSearchParams(filters).toString();
  const url = `${BASE_API_URL}/sessions?${queryParams}`;

  return fetch(url, {
      method: "GET",
      headers: {
          "Content-Type": "application/json",
          "SessionID": sessionId,
      },
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Failed to fetch sessions");
      }
      return response.json();
  })
  .catch(error => {
      console.error("Error fetching sessions:", error);
      return [];
  });
}



// Function to send the current user's location to the server
export const sendLocation = () => {
    const sessionId = localStorage.getItem("sessionid");

    if (!sessionId) {
        console.error("Session ID is missing. Redirecting to login...");
        window.location.href = "/login.html";
        return;
    }

    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
            (position) => {
                console.log("Sending coordinates to the server...");
                const locationData = {
                    sessionid: sessionId,
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude,
                };

                // Send location to the backend
                fetch(`${BASE_API_URL}/location`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(locationData),
                })
                .then(response => {
                    if (response.ok) {
                        document.getElementById("locationStatus").innerText = "Location sent successfully!";
                        console.log("Location sent successfully!");

                        // Initialize or center the map at the current user's location
                        initMap(position.coords.latitude, position.coords.longitude);
                    } else {
                        document.getElementById("locationStatus").innerText = "Failed to send location.";
                        console.error("Failed to send location:", response.statusText);
                    }
                })
                .catch(err => console.error("Error sending location:", err));
            },
            (error) => {
                document.getElementById("locationStatus").innerText = "Error retrieving location.";
                console.error("Error retrieving location:", error);
            }
        );
    } else {
        document.getElementById("locationStatus").innerText = "Geolocation not supported.";
        console.error("Geolocation is not supported by this browser.");
    }
};

// Function to fetch all user locations from the server
function fetchUserLocations() {
  const sessionId = localStorage.getItem("sessionid"); 
  
  fetch(`${BASE_API_URL}/location`, {
      method: "GET",
      headers: {
          "Content-Type": "application/json",
          "SessionID": sessionId,
      },
  })
  .then(response => {
      if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
  })
  .then(data => {
      console.log("User locations received:", data); // Inspect the response structure
      if (!data.locations) {
          console.warn("No locations available.");
          return;
      }
      console.log("User locations received:", data.locations);
      addMarkers(data.locations); // Pass the data to addMarkers
  })
  .catch(error => {
      console.error("Error fetching user locations:", error);
  });
}

