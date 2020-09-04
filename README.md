# Cafebucks

Simulating a coffee shop in go using microservices architecture. This is the repository for common configurations for all the services, like interacting with kafka etc.

## Plan

At first, i am going to make http endpoints for all the services returning json creating events and pushing events on kafka, then eventually other than gateway all the services are going to be communicating via grpc using protocol buffers. And in between i am planning to integrate all the CI/CD things. I am going to use Github Actions for this. Planning to integrate Circuit Breaker, Bulkhead Patterns for now. Then will try to implement CQRS.

## Where am i on the plan

Just starting for now. Going to add 3 services.

* Order Service (which is the starting point for now.)

* Beans Service (which will validate if the cafe has enough beans to complete the order.)

* Brewerie Service (which will listen to Order Accepted event and start brewing.)

## What is it for

Just for learning purpose.

## Links to services

* https://www.github.com/vishal1132/cafebucks-beans
* https://www.github.com/vishal1132/cafebucks-brewerie
* https://www.github.com/vishal1132/cafebucks-order
