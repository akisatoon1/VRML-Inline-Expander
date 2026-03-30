# VRML Inline Expander

[日本語版はこちら](README_ja.md)

A CLI tool that expands `Inline` nodes in VRML files and embeds the referenced file contents as `Group` nodes.

## Overview

This tool automatically detects `Inline` nodes in VRML 2.0 format files, expands the referenced file contents, and merges them into a single file. It's useful when you want to consolidate scenes split across multiple VRML files into a single file.

## Usage

### Installation

```bash
go build -o vrml-inline-expander
```

### Execution

```bash
./vrml-inline-expander <input-file> <output-file>
```

**Parameters:**
- `<input-file>`: Path to the VRML file to be expanded
- `<output-file>`: Path where the expanded VRML file will be saved

### Example

Suppose you have the following VRML files:

**top.wrl (input file):**
```vrml
#VRML V2.0 utf8

Inline {
  url "refered.wrl"
}
```

**refered.wrl (referenced file):**
```vrml
#VRML V2.0 utf8

Shape {
  geometry Sphere {}
}
```

Run the following command to expand:

```bash
./vrml-inline-expander top.wrl merged.wrl
```

**merged.wrl (output file):**
```vrml
#VRML V2.0 utf8
Group { children [Shape { geometry Sphere { } } ] }
```

The `Inline` node has been replaced with a `Group` node, and the referenced file contents have been embedded.

## Supported Formats

- **VRML Version**: VRML 2.0
- **Character Encoding**: UTF-8
- **Reference Type**: Local file paths only (network references like HTTP are not supported)

## Limitations

- The `url` field of `Inline` nodes must be in single string format (list format is not supported)
- Duplicate references to the same file are not expected
- Assumes no DEF name collisions between files

## System Configuration
[system configuration](/system.md)

## License

MIT License
