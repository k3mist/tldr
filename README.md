# tldr in golang

[TLDR pages](https://tldr-pages.github.io/) - Simplified and community-driven man pages

![Terminal](terminal.png)

## Install

```
go install bitbucket.org/djr2/tldr@latest
```

## Download

* [Distributions](https://bitbucket.org/djr2/tldr/src/master/dist/)

## Building and Build Requirements

* go
* upx -- https://github.com/upx/upx
* tar

```
bin/build
```

The build script will compile and compress the tldr executables.

The build script currently supports the following platforms and architectures;

* darwin arm64
* darwin amd64
* linux arm64
* linux amd64
* windows amd64

To a build a specific platform version run the below commands. 
It is important to replace `[platform]` with the desired operating system and 
`[arch]` with the desired platform architecture to build the executable correctly.

Supported Go build platforms and architectures can be found here;
https://golang.org/doc/install/source#environment

It is not necessary to run upx but it greatly reduces executable size.

```bash
GOOS=[platform] GOARCH=[arch] go build -ldflags="-s -w" -o tldr
upx --brute tldr # executable compression
```

## Usage

```
Usage:
  -L, --language string
        The desired language for the tldr page.
  -c page
        Clear cache for a specific tldr page.
        -p is required if clearing cache for a specific platform.
  --clear
        Clear the entire page cache.
  --debug string
        Enables debug logging. (default "disable")
  -g, --get
        If a tldr page is not cached attempt to retrieve it.
  --help
        This usage output.
  -p, --platform string
        Platform of the desired tldr page.
  --platforms
        Display a list of available platforms.
  -u, --update
        Update the local page cache.
  --version
        Display the version number.
```

### View a tldr
```
tldr <page>
```

### View a tldr for a specific platform
```
tldr <page> -p osx
```

### View a tldr for a specific language
```
tldr <page> -L en
```

### View a tldr for a specific language and platform
```
tldr <page> -L en -p osx
```

### Clear a tldr
```
tldr -c <page>
```

### Clear a tldr for a specific platform
```
tldr -c <page> -p osx
```

### Clear a tldr for a specific language and platform
```
tldr -c <page> -L en -p osx
```

### Clear entire cache
```
tldr --clear

or

tldr -c clearall
```

### Usage Notes

CLI arguments may be provided in any order.

For example;
```
tldr -p windows cd -L it
tldr -L es tar
```

## Configuration

A configuration is created the first time `tldr` is run.

The configuration is located at;
```
$HOME/.tldr/config.json
```

Pages repository URI, Zip URI, and all of the output colors are
configurable.

Below is the default configuration.

```
{
"pages_uri": "",
"zip_uri": "",
"language": "",
"cache_expiration": 30,
"extended_search": true,
"lookup_warnings": false,
"banner_color_1": 36,
"banner_color_2": 34,
"tldr_color": 97,
"header_color": 34,
"header_decor_color": 97,
"platform_color": 90,
"platform_alt_color": 95,
"description_color": 0,
"example_color": 36,
"hyphen_color": 0,
"syntax_color": 31,
"variable_color": 0
}
```

If plain (default) terminal text is desired set all color options to `0`.

`pages_uri` and `zip_uri` when left blank will use the official TLDR
locations.

These can be used to test pages from a custom repository
or any zip collection that follows the official TLDR directory format
and file specification.

Pages: `https://raw.githubusercontent.com/tldr-pages/tldr/main/pages/`

Zip: `https://tldr-pages.github.io/assets/tldr.zip`

To reset the configuration back to its defaults delete `config.json`
and it will be recreated. Or copy and paste the configuration from
this README above.

## License

[MIT License](https://bitbucket.org/djr2/tldr/src/master/LICENSE.md)
