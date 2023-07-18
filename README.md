# Shopping list telegram bot


## TODOs:
- Add concurrency to the handlers
- Research into how the web hooks will work and if there's any concurrency built in with them
    - Aka, each update comes in into its own goroutine 
- Add memcashed or redis and use it for timed cleanup
- Need to figure out if i can make one db transaction per journey rather then multiple, maybe using the in mem db?
- Need to figure out how to break out of journeys before they are finished, this will even help with the infinite loop




