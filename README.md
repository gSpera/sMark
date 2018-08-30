sMark {#sMark}
-----

sMark is a simple markup language built for humans, it is easy to read
and write,
you can get an introduction [here](https://gSpera.github.com/sMark)

Installing {#Installing}
----------

If you already have Golang installed you can get it throught go get
(go get https://github.com/gSpera/sMark)

Otherwise you can downloa9d the prebuilt binary from
[Github](https://github.com/gSpera/sMark/releases)

Usage {#Usage}
-----

Use sMark - help to get a brief intruduction to all command parameters\

Basicaly you could use
sMark -i inputfile.sm -o output.html

The main output engine produces HTML but there are other output engines
and you can make you own

Telegra.ph: sMark -i inputfile.sm -telegraph
Prettifier: sMark -i inputfile.sm -prettify = output.sm

WIP {#WIP}
---

sMark is in a developing state, some functions are now documented and
the code is not optimal

README.sm {#README.sm}
---------

The file you are viewing is a Markdown version of the file README.sm,
there is no easy way to convert sMark to Markdown(you can write you own Output Engine), it is actualy
compiled to HTML and then converted to Markdown
This solution is not optimal but GitHub doesn\'t support sMark
