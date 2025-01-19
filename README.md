# natto

austere gemini/spartan tools for openbsd (might work on other platforms...?)

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
protocol "spartan"

relay "gemini" {
  listen on 0.0.0.0 port 1965 tls
  protocol gemini
  forward to ::1 port 1965
}

relay "spartan" {
  listen on 0.0.0.0 port 300
  protocol spartan
  forward to ::1 port 300
}
```

and inetd.conf:
```
[::1]:1965 stream tcp6 nowait gemini /usr/local/bin/natto natto
[::1]:300 stream tcp6 nowait gemini /usr/local/bin/natto natto -s
```

## karashi

standalone gemini server. handles tls.

## negi

standalone spartan and gemini server. doesn't handle tls.

## okra

very simple gemini client, mostly for testing purposes

