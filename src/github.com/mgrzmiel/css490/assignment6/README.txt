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
To run the program, you have to type: go run timeserver.go (plus arguments) and in another window go authserver.go (plus arguments)
Arguments:
	--port <portNo>     - (optional) defines port on which server is running (default: 8080)
	-V                  - (optional) outputs version of the server and terminates it
	-—templates         - (optional) to provide the path for templates (default: “templates/“)
	—-log               - (optional) to specify the name of the log configuration file (default: “etc/timeserverlog.xml”)
	—-authport <portNo> - (optional) to define the port on which the authserver is running (default: 8070)
	—-authhost	 	    - (optional) to specifying the hostname of the authserver (default: localhost)
	—-authtimeout-ms    - (optional) the time when the timeserver is waiting to get the respond from authserver in milliseconds (default: 0)
	—-avg-response-ms   - (optional) the average response time to generate the delay in respond to the time request (default: 0)
	—-deviation-ms      - (optional) standard deviation to generate the delay in response to the time request (default: 0)
	—-max-inflight      - (optional) the maximum number of requests that the server can handle, if not specified it handle as much as it  can
	—-dumpfile          - (optional) the name of the file from which the dictionary will be loaded and save every given period of time, 
                              if not specified no data will be saved
	—-checkpoint-interval - (optional) specifies time in milliseconds, every checkpoint-interval the data from dictionary is saved to the dumpfile, if not specified, data will not be saved 

To run the program with the load generator, you have to additionally run ./bin/loadgen with the following parameters:
	--url 				- URL to sample e.g.'http://localhost:portNo/time' 
	--runtime 			- number of seconds to process
	--rate 				- average rate of requests (per second)
	--burst				- number of concurrent requests to issue
	--time-out-ms		- max time to wait for response
	—-log               - (optional) to specify the name of the log configuration file (default: “etc/loadgenlog.xml”)


You can also used the Makefile to run the program - just type make in the command line and then bin/timeserver (plus arguments) and 
in another window bin/authserver (plus arguments) and bin/load (plus argument). You can also run all of the programs in one window 
(e.g. ./bin/authserver --log=etc/authserverlog.xml & ./bin/timeserver --log=etc/timserverlog.xml --port=8081 --max-inflight=80 --avg-response-ms=500 --response-deviation-ms=300 & ./bin/loadgen --log=etc/loadgenlog.xml --url='http://localhost:8081/time' --runtime=10 --rate=200 --burst=20 --timeout-ms=1000).


COPYING / LICENSE:	
Copyright 2015 Magdalena Grzmiel
