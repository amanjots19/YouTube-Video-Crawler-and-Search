# Fampay Assignment

A RESTful API to fetch the latest videos from YouTube for a given search query and store them in a database. The API allows users to access the stored videos in a paginated response sorted in descending order of their publishing date-time.

## Features
- Continuously calls the YouTube API in the background to fetch the latest videos for a defined search query
- Stores the video data in a database with proper indexes
- Provides a GET API to access the stored videos in a paginated response sorted by publishing date-time
- Implements a basic search API to search the stored videos using their title and description
- Dockerized for easy deployment and scaling
## Requirements
- Golang 1.15 or later
- Docker and Docker Compose
- YouTube API key
 
## Current Predefined Query
### Tea 
## Setup
### Clone the repository:
```shell
$ git clone https://github.com/[your_github_username]/assignment-fampay.git`

$ cd assignment-fampay

$ docker-compose up
```

## Also you can run the application locally using the following commands
```make 
$ make clean`

$ make build

$ make run
```



The API will be available at http://localhost:8080/

## Get API Get Videos

`http://localhost:8080/video?limit=10&offset=0`

Endpoint will respond with latest videos data as a paginated response

## Get API Search Videos

`http://localhost:8080/search?q=make tea&limit=10&offset=0`

Endpoint will respond with latest videos with query(q) and get a paginated response,
Query has been optimized so that user can search in any form like "how to make a tea" can be "tea how".

## License

[MIT](https://choosealicense.com/licenses/mit/)
