
# mongo 3.2 docker image based on basic ubuntu 15.10 64 bit

see docker location on docker hub
- https://hub.docker.com/r/jumpscale/mongo/

std system passwd
- root/gig1234

no authentication for mongodb is needed

## use from non jumpscale enabled system

easiest way to use (windows & max)
- install https://www.docker.com/toolbox
- login using your docker account (create one if you don't have yet)
- click create
- look for jumpscale/mongo

careful the std passwd is now used, use jsdocker for more security (see below)

## use from jumpscale enabled system

```
docker pull jumpscale/mongo
jsdocker new -n kdstest -b jumpscale/mongo --ports 9000:9000 --start
```

you can now login with
```
ssh localhost -p 9022
```
port will change depending nr of dockers on your machine

# to build

```
mkdir -p /opt/code/github/jumpscale
cd /opt/code/github/jumpscale
git clone https://github.com/Jumpscale/docker_mongo.git
#or
#git clone git@github.com:Jumpscale/docker_mongo.git
sh build_docker.sh
```

# remarks

## runit

this docker uses runit to schedule all processes some useful commands
- more info on http://smarden.org/runit/faq.html

imagine the docker name is 04c9611bd06a
```
#find docker id 04...
docker ps
#restart ssh
docker exec 04c9611bd06a sv restart sshd
```

## ssh troubles

ssh tries to use all private local keys
it could be your ssh server is configured to not allow too many attempts
to make sure that no ssh keys are used and only try with login/passwd combination do following
```
ssh -o PubkeyAuthentication=no localhost -p 9026
```
port you can find by doing 'docker ps'
