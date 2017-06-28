<p align="center">
<img
    src="logo.png"
    width="260" height="80" border="0" alt="cowyodel">
<br>
<a href="https://travis-ci.org/schollz/cowyodel"><img src="https://img.shields.io/travis/schollz/cowyodel.svg?style=flat-square" alt="Build Status"></a>
<a href="https://github.com/schollz/cowyodel/releases/latest"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg?style=flat-square" alt="Version"></a>
</p>

<p align="center">CLI tool for interacting with <a href="https://github.com/schollz/cowyo">cowyo</a>  :cow: :speech_balloon:</p>

*cowyo* is a self-contained wiki server that makes jotting notes easy and _fast_. The most important feature here is _simplicity_. Other features include versioning, page locking, self-destructing messages, encryption, and listifying. You can [download *cowyo* as a single executable](https://github.com/schollz/cowyo/releases/latest) or install it with Go. Try it out at https://cowyo.com.

Getting Started
===============

## Install

If you have go

```
go get -u -v github.com/schollz/cowyodel
```

or just download from the [latest releases](https://github.com/schollz/cowyodel/releases/latest).

## Run

```
cowyodel upload FILE
```

OR

```
cat FILE | cowyodel upload
```

## License

MIT