<HERE I WILL WRITE WHAT IS IN MY MIND AS I AM NOT USING AI IN THIS PROJECT>

First I need to have these

- a master node where i can register the Pi nodes (master)
- a daemon that i can install on Pi to connect it to master (agent)
- a web view for seeing all the stuff.

### Agent
This is not going to be a simple http server instead this will be a daemon.
I need to make it in a way that first when install it checks for all the deps and install them setup Raspberry Pi
for this use so that any one can install it in one click and not look back again.



think we should start my making master first as that will have apis where agent will register making agent first makes no sense as if installing and setting to things will be in script as our main problems are will cgnet and monitoring that will be controlled by master we need tailscale setup at master first with apis so we can use that to build agent