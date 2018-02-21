Sketchboard Tool
================

The official, but experimental [Sketchboard](https://sketchboard.io) tool.

Supports converting a CSV column to note shapes.


Installation
------------

Download binary from the releases or build Sketchboard tool from the command line.


Building
--------

go build -o sketchboard


Usage
-----

Copies note_shapes.csv first column to the clipboard (MacOSX).

```
sketchboard convertcsv -f test-data/note_shapes.csv | pbcopy
```

Change default note size
```
sketchboard convertcsv -f test-data/note_shapes.csv --note-width 250 --note-height 150
```

Paste Cmd+P (Ctrl+P) to your Skethboard board.


Limitations
-----------

Uses fixed note shape dimensions and note shapes are not automatically resized.
