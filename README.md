An attempt to implement a MapleStory server emulator in Go.
This project is still extremely early, so it only does very few things at the moment.
Most of this is put together by studying TitanMS, OdinMS and Vana so huge 
credits to them for figuring packet structures and other stuff out, this is merely a Go implementation.

Feature progress
============
* Initial handshake with the client - [done](http://www.hnng.moe/f/49)
* Packet decryption/encryption - [done](http://hnng.moe/f/4m)
* Login - [done (video)](http://hnng.moe/f/5N)
* MySQL - [done](http://www.hnng.moe/f/5H)
* World / Channel Selection - [done (video)](http://hnng.moe/f/6N)
* Character Selection / Creation / Deletion - [done](http://hnng.moe/f/7l)
* Handling Bans - [done](http://www.hnng.moe/f/5Q)
* Getting in game - WIP

Getting started
============
First of all, you need to install [maplelib](https://github.com/Francesco149/maplelib) and  [mymysql](https://github.com/ziutek/mymysql).

Once you've done that, all that's left is to clone the repository.
Make sure that you have git and go installed and run

    go get github.com/Francesco149/kagami


You can also manually clone the repository anywhere you want by running

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
You can get the documentation with the built-in godoc 

    godoc github.com/Francesco149/kagami
    
If you're looking for a specific function or type just use

    godoc github.com/Francesco149/kagami MyFunction