README:	
This is an example of http server which have implemented following endpoints:
- index.html - if user is not log in - then login form is displayed
	     - if user is log in them the greeting message is displayed
- login - if the name in query string is provided it redirects the page to the index.html endpoint and displays the greeting message. 
        - if the name is not provided, it displayed the following message “C'mon, I need a name.”
- logout - logout the user and display goodbye message
- time - if the user is log in display the time plus user name
       - if the user is not login displays the time.
For any other endpoints, it displays information "These are not the URLs you're looking for" and sets the status code to 404 - Not Found.

AUTHORS: 
Magdalena Grzmiel

RUN:	
To run the program, you have to type: go run timeserver.go
Arguments:
	--port <portNo>   - (optional) defines port on which server is running (default: 8080)
	-V                - (optional) outputs version of the server and terminates it
	-—templates       - (optional) to provide the path for templates (default: “templates/“)
	—-log             - (optional) to specify the name of the log configuration file (default: “etc/seemlog.xml”)

You can also used the Makefile to run the program - just type make in the command line and then bin/assignment3 (plus arguments)

COPYING / LICENSE:	
Copyright 2015 Magdalena Grzmiel
