Server Goodies
==============

***Under Construction*** :
- will soon be split into directory names that make more sense.
- adding dockerfile

Server Goodies are the odds n' ends that are needed to glue everything together
so that it actually runs successfully.  As the complexity of the goodies
increases, they should be converted away from scripts and rewritten as
full programs.



Updating the Server with systemctl
-------------------

### When Server Code has changed

To fetch new changes to the repo, re-install the server software, and then
restart the servers:

	go get -u github.com/tilegame/gameserver
	sudo systemctl restart tilegameserver



### When tilegameserver.service has changed

When the .service module for systemctl is changed, do this instead:

	go get -u github.com/tilegame/gameserver
	cd ~/go/src/github.com/tilegame/gameserver/goodies
	chmod u=rx UpdateService.sh
	./UpdateService.sh

The UpdateService script will move the updated .service file into
the right location, and restart the daemon.





Using Docker
--------------------------

	TODO
	- write a VM startup script.
	- update server code with docker
	- figure out TLS when using docker.

Alternative approach to starting the game server.
Should make things easier to maintain and update,
and easier to add different parts (like databases).

There may be a need to do the `systemctl` for running docker.


### 1. Download and Install Docker from Terminal

There's an easy install script at https://github.com/docker/docker-install

	curl -fsSL https://get.docker.com -o get-docker.sh
	sh get-docker.sh

If Docker installation is successful, this warning is printed:

> If you would like to use Docker as a non-root user, you should now consider
adding your user to the "docker" group with something like:
>
> 	sudo usermod -aG docker fractalbach
>
> Remember that you will have to log out and back in for this to take effect!

To add the current user to the 'docker' group:

	sudo usermod -aG docker $USER

See https://docs.docker.com/install/linux/linux-postinstall/
for more information about what to do after installing docker.


### 2. Start Docker on Boot

More info:  https://docs.docker.com/install/linux/linux-postinstall/#configure-docker-to-start-on-boot

	sudo systemctl enable docker


### 3. Get the repo code

	TODO: docker pull something




Future Goals
------------

* be able to initialize a fresh VM automatically

* be portable to different environments (so that you don't always need a VM)
	* docker, perhaps?
	* various cloud platforms

* run the service in its own user group
	* create user group if it doesn't exist.
	* use the existing user if it does.

* rewrite bash scripts into golang so there is a more unified "tile game server tool".



Ideally, this should be entirely automated, and should minimize the
usage of bash scripts as much as possible.  Updating a specific branch
of the repo should trigger the update process.

One possibility is to a program that handles the graceful startup and shutdown
of the various components of the tile game server.  At the time of this writing,
there aren't enough components in existence for this to matter.
But when there are, it would be wise to have something able to manage them.
