# G3 Hack Challenge at .devHacks 2024!

The following web application was initially created for [.devHacks](https://devclub.ca/devhacks). 

## Title: User Session and Location Monitoring Dashboard with Map View
Create a web application that monitors user sessions and locations in real-time, logging the data to a SQL server and displaying it in an interactive dashboard with a map view. When a user logs in to the web page, track their session duration and location and visualize this data on a map within the dashboard.

## Key Components:
### [User Authentication](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#user-authentication) 
Implement a secure user authentication system for accessing the web application.

### [Session Tracking](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#session-tracking)
Develop functionality to track user sessions, recording the start time when a user logs in and the end time when they log out or close the tab/browser.

### [Location Tracking](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#location-tracking)
Utilize browser geolocation API or IP-based geolocation to determine the user's location accurately. Log this information along with the session data to the SQL server.

### [SQL Database Integration](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#sql-database)
Set up a SQL database to store user session and location data efficiently. Design appropriate tables to maintain this data.

### [Real-time Updates](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#real-time-updates)
Implement real-time updates to the dashboard to reflect the latest user session and location data as it is logged to the SQL server. Utilize WebSockets or server-sent events for real-time communication.

### [Dashboard Visualization](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#dashboard-visualization-and-ui)
Design an interactive dashboard with a map view to visualize user location data. Use a mapping library like Leaflet.js or Google Maps API to display user locations dynamically on the map.

### [User Interface (UI)](https://github.com/niraaz1912/Codezilla-Repository/blob/main/details.md#dashboard-visualization-and-ui)
Create a user-friendly UI for the dashboard, allowing administrators to view and analyze user session and location data easily. Include features for filtering, searching, and customizing the displayed data.

### Performance Optimization
Optimize the performance of both the web application and the SQL server to handle a potentially large volume of user session and location data efficiently.

### Security Measures
Implement security measures to protect user data, including encryption of sensitive information and preventing unauthorized access to the dashboard.
