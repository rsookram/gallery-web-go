# gallery-web-go

`gallery-web-go` is a web app which displays the images within a directory
(including subdirectories). A Go server serves the page and images with one
image being displayed at a time. No third-party dependencies are used.

## Building

Build the server by running the following command in the root of this repo:

```shell
go build .
```

Or run the following to build and install:

```shell
go install .
```

## Usage

```
USAGE:
    gallery-web [DIR] [OPTIONS]

Start a server for the specified directory (defaulting to the current
directory).

OPTIONS:
    -port <port>  The port to serve on [default: 4255]
```

Open `http://localhost:4255/` in a browser to view the images.

## Key bindings

|Key|Description|
|---|-----------|
|Left arrow|Move to next image|
|Right arrow|Move to previous image|
|Up arrow|Zoom in|
|Down arrow|Zoom out|
|e|Toggle display of position in gallery|
