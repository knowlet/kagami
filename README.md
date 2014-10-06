An attempt to implement a MapleStory server emulator in Go.
This project is still extremely early, so it only does very few things at the moment.
Most of this is put together by studying TitanMS, OdinMS and Vana so huge 
credits to them for figuring packet structures and other stuff out, this is merely a Go implementation.

Support me!
============
Like my releases? Donate me a coffe!

Paypal: [click](http://hnng.moe/6M)

Litecoin: LUZm98D1nPhNQBw9QjkSS9XJee9X5hPjw3

Bitcoin: [15Jz8stcnkorzwCbUNk3qQbg2H9eySKXtb](bitcoin:15Jz8stcnkorzwCbUNk3qQbg2H9eySKXtb?label=donations) or [Bitcoin QR Code](http://hnng.moe/f/CM)

Dogecoin: [DDaYKDUxib2SnVEk9trG98ajCyu1Hw9zgQ](dogecoin:DDaYKDUxib2SnVEk9trG98ajCyu1Hw9zgQ?label=donations&message=wow%20much%20donate%20very%20thanks) or [Dogecoin QR Code](http://hnng.moe/f/CL)

Feature progress
============
* Initial handshake with the client - [done](http://www.hnng.moe/f/49)
* Packet decryption/encryption - [done](http://hnng.moe/f/4m)
* Login - [done (video)](http://hnng.moe/f/5N)
* MySQL - [done](http://www.hnng.moe/f/5H)
* World / Channel Selection - [done (video)](http://hnng.moe/f/6N)
* Character selection / creation / deletion and multiple worlds/channels - [done (video)](http://www.hnng.moe/f/7m)
* Handling Bans - [done](http://www.hnng.moe/f/5Q)
* Getting in game - [done](http://hnng.moe/f/Ai)
* Basic in-game sync - WIP ([video of portals working](http://hnng.moe/f/CK))
* Properly syncing data between login/world/chan - WIP
* Graceful server shutdown that logs all players off - WIP

Getting started
============
Make sure that you have git and go installed and run the following commands to acquire all of the requires libraries.

	go get github.com/jteeuwen/go-pkg-xmlx 
	go get github.com/Francesco149/maplelib
	go get github.com/ziutek/mymysql/thrsafe
	go get github.com/ziutek/mymysql/autorc
	go get github.com/ziutek/mymysql/godrv

You can test these libraries before building kagami if you want.
First of all, create the test mysql user and database from your mysql console:

	mysql> create database test;
	mysql> grant all privileges on test.* to testuser@localhost;
	mysql> set password for testuser@localhost = password("TestPasswd9");
	
Make sure that max_allowed_packet is set to at least 34M in your my.ini/my.cnf, then run the tests for the mymysql library:

	go test github.com/ziutek/mymysql/...
	
Now you can go ahead and test maplelib and other libraries:
	
	go test github.com/jteeuwen/go-pkg-xmlx/...
	go test github.com/Francesco149/maplelib/...

Once you've made sure that all of the libraries are working properly, you can obtain the actual server:

	go get github.com/Francesco149/kagami

If you want, you can also manually clone the repository anywhere you want by running

	git clone https://github.com/Francesco149/kagami.git
    
Before you run the server you will also need to configure your MySQL database 
info in kagami/common/consts/consts.go . Don't worry, this is temporary - 
everything will be moved to config files as soon as I get more things working.

Make sure that your MySQL database is running and make sure that you've created 
the kagami database by running the query in the kagami.sql file.

NOTE: the database structure will change very often at the current stage of the project and you might end up having to delete and recreate your database after an update.
    
Running the server
============
To compile the server, all you have to do is:

	go install github.com/Francesco149/kagami/...

And then simply run loginserver, worldserver and as many channels servers as you like in your $GOPATH/bin directory.
    
Documentation
============
You can view the documentation as HTML by simply running

	godoc -http=":6060"

and visiting

	http://localhost:6060/pkg/github.com/Francesco149/kagami/
