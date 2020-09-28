# Introduction
RipTrace is a microservice debugger designed to debug applications written in Golang. The debugger relies on Delve (ptrace) to do arbitrary debugging, and as such is not designed
for "live" production systems. Instead, RipTrace is designed to be used in conjunction with Docker Compose, Minikube, or similar to debug multiple microservices at the same time as they interact
in a networked environment. 

# NATS
RipTrace uses NATS as the primary orchestration and communication mechansim. Each debugger agent runs the following topics:
- "<host>.profile.get": Returns a profile of the host including all running go processes. The PID of the process is used to attach the debugger agent to the process.
- "<host>.debugger.attach": Given a PID returned from the profile.get topic, this will instruct the agent to attach a delve process to a running process.
- "<host>.debugger.createBreakpoint": Given a breakpoint location, create a non-breaking breakpoint and a new publisher topic. Every time the breakpoint is hit, RipTrace
   will emit a message on the "<host>.debugger.trace" topic containing the full trace information (locals, variables, line, file, etc)
   

# TODO:
- Add host information (most likely IP address) to the topic to be able to coordinate multiple hosts
- Consume breakpoint initialization file or similar (eventually will be a frontend IDE style with breakpoints)
