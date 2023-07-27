# Shopping list telegram bot
## Goal
- Be able to put this bot in a chat and have it track shopping lists
- With multiple people having access to this list in mind

### Use case
- One person always enters stuff in a shopping list and another accesses the list when shopping comes


## Wider goal
To learn and possibly create a framework for this telegram bot library in which I could go on to creating a personal telegram bot with modules which could be loaded.
Module example: 
- Server log management
- Server script management and execution
- Media management
    - For example be able to forward it songs on telegram then have the bot take the songs and place them in an appropriate folder on your personal NAS
- Get updated about services running on your server etc


## Shopping List Bot TODOs:
- [] Figure out how to make webhooks work with railway.app ssl
- [x] Add concurrency to the handlers
- [x] Need to figure out how to break out of journeys before they are finished, this will even help with the infinite loop
- [x] Need to replace built in logger with smth else
    - Replaced it with log15 but will probs switch to logrus
    - With log15 I cannot seem to find a way to set Info to output to stdout and separately Error to stderr
- [] Have a think about changing the handlers to have my own context which would then hold data on journeys, update, and possibly be able to mess about with the journey itself.
    - Example: be able to set the iterator for the journey index so you could use this the infinite loop, break out the journey, skip a journey index or go to the beginning
    - This context could also hold instances or logger which would be initialised already with chatID and other data
    - Might be a good idea to also get rid of the handler contextual return and put that in the context itself
- [] Clean up the logging in the entire app
    - I.E. there is too much logging going on
    - Think about making a debug level of logging
- [] Come up with a way to structure keyboards (inline as well) so that each button can have its own handler so I don't need to rely on ifs and switch cases statements
