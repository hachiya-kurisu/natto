# natto

austere gemini/spartan server for openbsd (might work on other platforms...?)

* inetd-based
* tls handled by relayd
* does unveil/pledge on openbsd
* spartan support 💪


## setup

relayd.conf will need something like this:
```
protocol "gemini" {
  tls keypair gemini
}

relay "gemini" {
  listen on 0.0.0.0 port 1965 tls
  protocol gemini
  forward to ::1 port 1965
}
```

and inetd.conf:
```
[::1]:gemini stream tcp6 nowait gemini /usr/local/bin/natto natto
```

you might have to define the gemini service in /etc/services:
```
gemini 1965/tcp
```

