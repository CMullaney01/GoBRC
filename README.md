# GoBRC

Completing the 1-Billion Row Challenge in Go!

# Attempt - 1
Lets start off basic, single threaded parse the file line by line and update our info for for each city as we go along.

achieved an impressive: 1m40.935s

# Attempt - 2
Two threads. Unlike our C++ attempt we have some sexy go channels to use I am going to use a two threads system, 1 parser and 1 reader

achieved: 1m41.520s

No improvement!?