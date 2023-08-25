# Mini-Twitter

A simple twitter backend RESTful service in Golang, powered by Gin

* User sign-up
    * POST `/users`
    * Request body
      ```json
      {
            "username": "Bradley",
            "password": "33521",
            "email": "bcooper@gmail.com"
      }
      ```
* Post tweet in text
    * POST `/tweets`
    * Request body
      ```json
      {
            "user": "Bradley",
            "text": "let's rock!"
      }
      ```
* Access timeline
    * GET `/timeline/username`

* Access tweets
    * GET `/tweets/username`
* Follow and Unfollow users
    * POST `users/follow` and `users/unfollow`
    * Request body
      ```json
      {
         "from": "Bradley",
         "to": "Duncan"
      }
      ```
