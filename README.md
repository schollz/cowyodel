<p align="center">
<img
    src="logo.png"
    width="260" height="80" border="0" alt="cowyodel">
<br>
<a href="https://travis-ci.org/schollz/cowyodel"><img src="https://img.shields.io/travis/schollz/cowyodel.svg?style=flat-square" alt="Build Status"></a>
<a href="https://github.com/schollz/cowyodel/releases/latest"><img src="https://img.shields.io/badge/version-1.2.3-brightgreen.svg?style=flat-square" alt="Version"></a>
<img src="https://img.shields.io/badge/coverage-43%25-yellow.svg?style=flat-square" alt="Code Coverage">
</p>

<p align="center">Easily move things between computers using  <a href="https://github.com/schollz/cowyo">cowyo</a>  :cow: :speech_balloon:</p>

*cowyodel* allows simple and secure sharing of text/data between computers.  *cowyodel* temporarily transfers your data (with optional client-side encryption) to [a cowyo server](https://github.com/schollz/cowyo) where it resides until the other computer downloads it using the provided secret code phrase.

Demo
====

[![asciicast](demo.gif)](https://asciinema.org/a/Oq6enXjipBXqFcugqV7mSvdpR)

Getting Started
===============

## Install

If you have Go1.7+

```
go get -u -v github.com/schollz/cowyodel
```

or just download from the [latest releases](https://github.com/schollz/cowyodel/releases/latest).

## Basic usage 

### Upload

To share a file with another computer, you first upload it to a cowyo server. By default *cowyodel* uses [cowyo.com](https://cowyo.com), but [you can host your own cowyo server](https://github.com/schollz/cowyo) as well (see Advanced Usage).

```
$ cowyodel upload README.md
Uploaded README.md (textual data). Your codephrase: total-perform-retro

View/edit your data:

	https://cowyo.com/total-perform-retro

Download using cowyodel:

	cowyodel download total-perform-retro
```

or

```
$ cat README.md | cowyodel upload
Uploaded (textual data). Your codephrase: total-perform-retro

View/edit your data:

	https://cowyo.com/total-perform-retro

Download using cowyodel:

	cowyodel download total-perform-retro
```

After uploading, you will recieve a code-phrase, in the above example the code-phrase is `total-perform-retro`. If you don't want to use code phrases, you can also specify your own name using `-name`, see Advanced Usage below.

The uploads are fully compatible with [the cowyo server](https://cowyo.com), so you can view and edit them using the code-phrase (e.g.  `cowyo.com/total-perform-retro` in above example). 

### Download

On any other computer connected to the internet, you can download the file using the name using `cowyodel download code-phrase`.

```
$ cowyodel download total-perform-retro
Wrote text of 'total-perform-retro' to 'README.md'
```

After downloading, it will be erased from the cowyo.com. If you don't trust this server, you can also specify your own (see Advanced Usage). You can add `--store` flag to prevent automatic deletion from the cowyo server.


Advanced Usage
===============

### Persist (don't delete after first access)

Adding `--store` will keep the file on the server, even after accessing.

```
$ cowyodel upload --store FILE
```


### Specify codephrase

If you don't want to use a codephrase, use `--name` to specify the name of file.

```
$ cowyodel upload --name README.md
Uploaded README.md (textual data). Your codephrase: README.md
```

It is possible that someone could have used that page (and locked it) which would not allow that page to be used and a message "Locked, must unlock first" will appear.

### Client-side encryption

```
$ cowyodel upload --encrypt README.md
Enter passphrase: 123
Uploaded README.md (textual data). Your codephrase: total-perform-retro

$ cowyodel download total-perform-retro
Enter passphrase: 123
Wrote text of 'total-perform-retro' to 'README.md'
```

The encryption is fully compatible with the server-side encryption on [cowyo.com](https://cowyo.com), so you can still use the web browser to decrypt/encrypt your document.

If the decryption fails, the document will be re-uploaded to the cowyo server.

### Binary files

*cowyodel* automatically detects binary data and uploads as such.

```
$ cowyodel upload /tmp/image.jpg
Uploaded image.jpg (binary data). Your codephrase: empire-cricket-tokyo

Download using cowyodel:

	cowyodel download empire-cricket-tokyo

$ cowyodel download image.jpg
Wrote binary data to 'image.jpg'

$ sha256sum /tmp/image.jpg image.jpg
62a9583758d54e666ff210be3805483bd76ac522ea649f0264de65124943c0b3 */tmp/image.jpg
62a9583758d54e666ff210be3805483bd76ac522ea649f0264de65124943c0b3 *image.jpg
```

_Note:_ you should not access uploaded binary files at via the web browser as it would risk corrupting them.

### Self-hosting cowyo server

You can also [host your own cowyo server](https://github.com/schollz/cowyo) and use that instead of the default `cowyo.com`. To host *cowyo* yourself, just use

```
$ go get github.com/schollz/cowyo/...
$ cowyo
Running cowyo server (version ) at http://localhost:8050
```

(If you don't have Go installed, [you can also download a release version](https://github.com/schollz/cowyo/releases/latest)).

Once you have a self-hosted cowyo server, you just need to specify the server when running *cowyodel*:

```
$ cowyodel --server http://localhost:8050 upload FILE
uploaded to total-perform-retro
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


Inspiration
===========

This tool was inspired by the following:

- [wormhole](https://github.com/warner/magic-wormhole)
- [piknik](https://github.com/jedisct1/piknik)

*cowyodel* does not represent a significant innovation over these tools. However, there are some advantages that *cowyodel* provides:

- Trust. You can run your own *cowyo* server on a domain you trust.
- Direct edting. You can directly edit plaintext documents on the *cowyo* server using the web interface.
- Simplicity. The codebase is < 1k LOC, and is straightforward to understand.

Development
===========

To run tests, make sure to start a `cowyo` server first.

```
$ go get -u github.com/schollz/cowyo/...
$ cowyo
$ go test
```

License
========

MIT