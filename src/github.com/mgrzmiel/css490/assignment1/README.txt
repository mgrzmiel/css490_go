
README:	
This is a simple example of http server which display the actual time
after receiving following request: /time.
Otherwise, it displays information "These are not the URLs you're looking for"
and sets the status code to 404 - Not Found.

AUTHORS: 
Magdalena Grzmiel

RUN:	
To run the program, you have to type: go run timeserver.go
Arguments:
	--port <portNo>   - (optional) defines port on which server is running (default 8080)
	-V                - (optional) outputs version of the server and terminates it

COPYING / LICENSE:	
Copyright 2015 Magdalena Grzmiel
