# Postback Delivery

This is a simple postback delivery system build with php, golang, and Redis.

## Getting Started

Requirements to run the project:
- minimum go 1.6
- minimum php 5.6.0
- curl
- Redis

## How To Use

* Start a local built in php server that requests can be made to.
* Send a few curl request to [ingest.php].
* Run [delivery.go] to start processing request (can continue running while adding more request to [ingest.php])
* Logs are written to a file called logs.txt in the go folder.

## Documentation

[ingest.php] accepts a JSON object of the format:

(POST) http://{server_ip}/ingest.php
    (RAW POST DATA) 
    {  
      "endpoint":{  
        "method":"GET",
        "url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"
      },
      "data":[  
        {  
          "mascot":"Gopher",
          "location":"https://blog.golang.org/gopher/gopher.png"
        }
      ]
    }


It then parses the JSON into strings and puts them in the Redis database as follows:

* Checks if an ID for this system is in the database, if not it creates one starting at 1000
* Puts the "endpoint" object into a Hash set called "endpoint:xxxx" where xxxx is the ID generated above.
* The "data" object is also put in a Hash set called "data:xxxx" where xxxx is the ID generated above.
* After both objects have been succesfully pushed the ID itself is pushed onto a stack to be pulled by the delivery agent.
* Lastly the ID is incremented for the next incoming request.

The objects are placed in their own separate "tables" to keep persistent and lightweight data.

[delivery.go] checks the Redis stack that [ingest.php] is pushing objects onto.

* Opens a file called logs.txt to write logs to.
* Establishes a connection to the Redis database.
* Pulls an ID from the "endpoint_query" stack and uses that ID to get the correct endpoint and data objects. If there are no ID's the program keeps looping and checking.
* The request is formatted with the data inside of it, if there are multiple sets of data for the same endpoint, several requests are created with each set of data.
* A basic client is created and the requests are sent.
* The delivery, response time, response code, and response body are all logged and timestamped.

The formation of the requests and handling of them are completed in a go routine per endpoint.


## Built With

* [phpredis] (https://github.com/phpredis/phpredis) - The php framework used to connect to Redis
* [go-redis] (https://github.com/go-redis/redis) - The golang framework used to connect to Redis

## Authors

* **Kaleb Matherne**


