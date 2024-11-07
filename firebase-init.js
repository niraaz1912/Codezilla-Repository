// firebase-init.js

// Firebase App (the core Firebase SDK)
import { initializeApp } from "https://www.gstatic.com/firebasejs/11.0.1/firebase-app.js";
import { getAnalytics } from "https://www.gstatic.com/firebasejs/11.0.1/firebase-analytics.js";

// Your Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyDDAjj0hWPeH1ekIdW8U1K6JU-sI1TZl_o",
  authDomain: "g3-hack-challenge.firebaseapp.com",
  projectId: "g3-hack-challenge",
  storageBucket: "g3-hack-challenge.firebasestorage.app",
  messagingSenderId: "640762127603",
  appId: "1:640762127603:web:21d8c8663b95be3bd7c386",
  measurementId: "G-8CB67KDDXB"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const analytics = getAnalytics(app);
