# speed-engine

We're going to build application about game.


**How To Run:**
1. You should install postgres on your docker:
    - docker pull postgres
    - docker run --name synapsis-postgres -e POSTGRES_PASSWORD=password -p 5434:5432 -d postgres
2. You should install redis on your docker:
    - docker pull redis
    - docker run --name synapsis-redis -p 6379:6379 -d redis