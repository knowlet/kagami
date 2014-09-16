An attempt to implement a MapleStory server emulator in Go.
This project is still extremely early, so it only does very few things at the moment.

Feature progress
============
* Initial handshake with the client - [done](http://www.hnng.moe/f/49)
* Packet decryption/encryption - [done](http://hnng.moe/f/4m)
* Login - WIP
* MySQL - WIP

Getting started
============
First of all, you need to install [maplelib](https://github.com/Francesco149/maplelib).

Once you've done that, all that's left is to clone the repository.
Make sure that you have git and go installed and run

    go get github.com/Francesco149/kagami


You can also manually clone the repository anywhere you want by running
    git clone https://github.com/Francesco149/kagami.git
    
Running the server
============
For now only the loginserver is present, so all you have to do is

    go install github.com/Francesco149/kagami/loginserver
    loginserver
    
Documentation
============
You can get the documentation with the built-in godoc 

    godoc github.com/Francesco149/kagami
    
If you're looking for a specific function or type just use
    godoc github.com/Francesco149/maplelib MyFunction