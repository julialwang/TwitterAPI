# TwitterAPI

Twitter feed API exercise that has the following endpoints:
* Create/delete user (/register & /delete)
* Login/Logout user (/login & /logout)
* Follow/Unfollow user (/follow & /unfollow)
* Post a tweet / delete a tweet (/tweet & /untweet)
* Get user info (/profile)
* Get following timeline (/timeline)

This project was created out of personal interest during my summer internship as a way to familiarize myself with Golang and the Docker / database environment relationships that I would have to manage moving forward in my project. I had creative liberty to choose how I wanted to go about doing this, and ended up settling on a Twitter mimic because of the variety of options for endpoints that I would be able to incorporate. 

I was in charge of all relevant design choices, such as my use of MongoDB rather than a relational database like sqlite, which was motivated solely by my interest in learning how to use a document-based database system. Each feature was tested with a variety of test cases through Postman, which I was able to familiarize myself with over the course of development. One of my coworkers participated in the testing process as well, where he gave me edge cases to demonstrate application robustness.
