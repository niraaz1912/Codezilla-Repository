# User Authentication

## Creating Account

When signing up, it receives JSON with username and password.
It queries the database if the username already exists, if it does, it returns 409 status.

![Screenshot 2024-11-07 040000](https://github.com/user-attachments/assets/a084d8f9-c46a-4183-9969-ef40475953f5)

It uses bcrypt to hash the password. Then inserts new user with hashed password into the database with default role.

![Screenshot 2024-11-07 113555](https://github.com/user-attachments/assets/b7eca0ef-d70d-449b-b7f3-329ab13405b1)

The following are the logs in the server during signing up:

![Screenshot 2024-11-07 035829](https://github.com/user-attachments/assets/6d95e99e-7272-47e1-9507-d7b1c1fe5544)


## Logging in

When logging in, it receives JSON with username and password. 

It checks if the username exists, if it does not, it returns 401 status. It then compares the hashed password using bcrypt and if successful, returns with 200.

![Screenshot 2024-11-07 040043](https://github.com/user-attachments/assets/795c776d-aff4-405b-a2a8-086b6c0a4a25)

The following logs are generated during successful login:

![Screenshot 2024-11-07 040126](https://github.com/user-attachments/assets/e5a930c0-c8d1-4249-84b9-a3e99b1c06aa)


## Invalid Login

Invalid logins are handled by returning 401 status.

![Screenshot 2024-11-07 040418](https://github.com/user-attachments/assets/66d6eb06-b9cb-4198-b83c-561ed84fdf42)

![Screenshot 2024-11-07 040454](https://github.com/user-attachments/assets/bcf2a218-f16c-48c3-a5b5-f6ded825d231)


# Session Tracking

Upon successful login, the program generates a new session ID for the user and then stores the start time of their session.

After logging out, it deletes the session ID and stores the session's end time.

![image](https://github.com/user-attachments/assets/da260d32-5ac2-4799-bd36-d1d5486ea0ee)


# Location Tracking

After the user logs in, they can send their location data to the server through a button, which requires browser permission. 

![image](https://github.com/user-attachments/assets/dcd2b872-5d01-42b5-abeb-4bfb40fd3700)

The location is logged in the SQL database.

![Screenshot 2024-12-24 000642](https://github.com/user-attachments/assets/c2746679-31e5-4f74-86d5-941585a470f2)


# SQL Database

There are 2 tables.
1. Users: username, role, passhash, sessionID, start_time, end_time
2. Location History: username, latitude, longitude, time


# Real time updates

All the data are updated in real time.


# Dashboard Visualization and UI

When anyone logs in with a admin role, they are redirected to dashboard where they get the insight of the active users, their location and session.
When hovering over the markers on the map, it displays the username.

![Screenshot 2024-12-24 000528](https://github.com/user-attachments/assets/6e55050d-594c-4fe1-9a9b-1e84d85bc8ae)



