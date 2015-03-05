# Warper

Simple tool for sharing files to your friend/collegues by simple curl command.

Example launch command is like this:

```$ sudo ./warper -port=80 -temp=/ -basedir=warper```

You can add alias for simple warp file to server:
```
warp() {if [[ -z $1 ]]; then echo "no file specified\n"; else file=$1; curl --upload-file $file "http://HOSTNAME/$file"; fi;}
alias warp=warp
```
and then just do something like this:

```$ warp somefile.png```


##### Docker (advanced way):
```
$ docker run -d -p 80:80 --name nginx -v /tmp/nginx:/etc/nginx/conf.d -t nginx

$ sudo cp templates/* /tmp/templates

$ docker run -d --name nginx-gen --volumes-from nginx \
  -v /var/run/docker.sock:/tmp/docker.sock \
  -v /tmp/templates:/etc/docker-gen/templates \
  -t jwilder/docker-gen:0.3.4 -notify-sighup nginx -watch --only-published /etc/docker-gen/templates/nginx.tmpl /etc/nginx/conf.d/default.conf

$ sudo mkdir -p /dock_vols/warper

$ docker run -d -e VIRTUAL_HOST=somehost.com -p :8080 -v /dock_vols/warper:/tmp/warper warper
```
