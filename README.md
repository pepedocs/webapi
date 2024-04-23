# webapi
A modular HTTP API server that supports websocket connections. 

# About
This is a submission for a certain golang assessment. The assessment is made to be ambiguous in terms of how it can be interpreted. For my case, I could've just written all these into one file and the same requirements will be met. 

However, I wanted to challenge myself by implementing a composable and extensible go HTTP API server so here this is the result. 

# Components

**cmd/server.go** - Assembles the top-level components of the application and implements a CLI.

**server** - Layer 7 servers concerned with implementations such as websocket and http/api servers.

**controllers** - Controller instances that are only concerned with dispatching client requests to services.

**controllers** - Service instances that are only concerned with providing specific services such as providing an interface to get the current swatch time.

**integration** - Integration tests.

  
