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

