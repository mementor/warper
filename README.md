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
