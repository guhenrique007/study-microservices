**Initialization**
* export dependecies and env vars
* initialize services ` go run .`
* (or docker way)

**Stack**
* Go
* Rabbit MQ

**Services**
* Match service
    * provides match results for matches.json  

* Championship service 
  * show match results in template:
  ```
    Palmeiras 4 X 0 Corinthians
    Santos 2 X 2 Atlético-MG
  ```

* Table 
    * process table by results and send to queue:
  ```
    Palmeiras — 3
    Atlético-MG — 1
    Santos — 1
    Corinthians — 0
  ```
