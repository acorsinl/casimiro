#Casimiro
A basic template for REST servers in golang created for personal use.

There is a lot of stuff to be done yet, 
check [TODO.md] (https://github.com/acorsinl/casimiro/blob/master/TODO.md)
for further information.

#Setup
##server.go
###Constants
Two environment variables are needed, one for the HTTP listening port, one for
the MySQL connection string. These two are set as constants so you can modify
their name as preferred.

Casimiro is supposed to run behind an API manager or similar proxy tools,
therefore it expects the user id to be given by the upper layer in a Header.
Name of that header can be changed in UserHeader constant.

ResourcesUrl: For each resource Casimiro defines a new file with all the 
standard REST methods, hence more constants like this should be added for 
each resource your server will serve. Names for the urls are set here.

###Routing
Just replace the variable names for your choice of preference.

##resources.go
Just a basic template for the basic REST methods. SQL queries must be added, as
well as extra validations needed for your business logic and extending/modifying
the Resource struct.

#Licensing
Casimiro is licensed under BSD 3 clause license. 
See [LICENSE] (https://github.com/acorsinl/casimiro/blob/master/LICENSE) for 
the full license text.