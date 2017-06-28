<p align="center">
<img
    src="logo.png"
    width="260" height="80" border="0" alt="cowyodel">
<br>
<a href="https://travis-ci.org/schollz/cowyodel"><img src="https://img.shields.io/travis/schollz/cowyodel.svg?style=flat-square" alt="Build Status"></a>
<a href="https://github.com/schollz/cowyodel/releases/latest"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg?style=flat-square" alt="Version"></a>
</p>

<p align="center">Easily move things between computers using  <a href="https://github.com/schollz/cowyo">cowyo</a>  :cow: :speech_balloon:</p>

*cowyodel* allows easy and simple sharing of text/data between computers.  *cowyodel* is a command-line tool that allows simple interaction with [a cowyo server](https://github.com/schollz/cowyo), providing a simple way to upload/download text/binary that is encrypted/unencrypted.

Demo
====

[![asciicast](https://asciinema.org/a/Oq6enXjipBXqFcugqV7mSvdpR.png)](https://asciinema.org/a/Oq6enXjipBXqFcugqV7mSvdpR)

Getting Started
===============

## Install

If you have Go1.7+

```
go get -u -v github.com/schollz/cowyodel
```

or just download from the [latest releases](https://github.com/schollz/cowyodel/releases/latest).

## Basic usage 

### Upload a document

To share a document with another computer, you first can upload it to the cowyo instance using `cowyodel upload FILE`.

```
$ cowyodel upload README.md
uploaded to 2-adoring-thompson
```
or
```
$ cat README.md | cowyodel upload
uploaded to 2-adoring-thompson
```

The uploads are fully compatible with [cowyo](https://cowyo.com), so you can reach them at the specified name (e.g.  `cowyo.com/2-adoring-thompson` in above example) to view/edit. You can also specify your own name using `-name`, see Advanced Usage below.

### Download the document

On any other computer connected to the internet, you can download the file using the name using `cowyodel download NAME`, where `NAME` is also the URL you can access it (e.g. cowyo.com/NAME).

```
$ cowyodel download 2-adoring-thompson
Wrote text to '2-adoring-thompson'
```

By default, the first time you access it (via web or downloading), it will be erased. To prevent this, you can add `--store`.


Advanced Usage
===============

### Persist (and don't delete after first access):

```
$ cowyodel upload --store FILE
```


### Specify filename

```
$ cowyodel upload --name README.md
uploaded to README.md
```

It is possible that someone could have used that page (and locked it) which would not allow that page to be used and a message "Locked, must unlock first" will appear.

### Client-side encryption

```
$ cowyodel upload --encrypt README.md
Enter passphrase: 123
uploaded to 2-adoring-thompson

$ cowyodel download 2-adoring-thompson
Enter passphrase: 123
wrote text to '2-adoring-thompson'
```

The encryption is fully compatible with th server-side encryption on [cowyo.com](https://cowyo.com), so you can still use the web browser to decrypt/encrypt your document.

### Binary files

Binary files are Gzipped and then Base64 encoded for transfering to/from the server. Thus, you should not access them via the web browser as it would risk corrupting them.

```
$ cowyodel upload --direct image.jpg
uploaded to 2-adoring-thompson

$ cowyodel download --direct image.jpg
wrote binary data to '2-adoring-thompson'

$ sha256sum image.jpg 2-adoring-thompson
62a9583758d54e666ff210be3805483bd76ac522ea649f0264de65124943c0b3 *logo.jpg
62a9583758d54e666ff210be3805483bd76ac522ea649f0264de65124943c0b3 *2-adoring-thompson
```

### Self-hosting files

You can also [host your own cowyo server](https://github.com/schollz/cowyo) and use that instead of the default `cowyo.com`. 

```
$ cowyodel --server myserver.com upload FILE
uploaded to 2-adoring-thompson
```

### Help

```
$ cowyodel -h
NAME:
   cowyodel - upload/download encrypted/unencrypted text/binary to cowyo.com

USAGE:
   Upload a file:
    cowyodel upload README.md
    cat README.md | cowyodel upload
   
   Download a file:
    cowyodel download 2-adoring-thompson

   Persist (and don't delete after first access):
    cowyodel upload --store FILE

   Specify filename:
    cowyodel upload --name README.md

   Client-side encryption:
    cowyodel upload --encrypt README.md

   Binary-file uploading/downloading:
    cowyodel upload --direct --name image.jpg
    cowyodel download --direct image.jpg

COMMANDS:
     upload, u    upload document
     download, d  download document
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server value  cowyo server to use (default: "https://cowyo.com")
   --debug         debug mode
   --help, -h      show help
   --version, -v   print the version
```

License
========

MIT
