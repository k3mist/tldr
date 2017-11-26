# tldr in golang

[TLDR pages](https://tldr.sh) - Simplified and community-driven man pages

![Terminal](terminal.png)

## Install

```
go get -u bitbucket.org/djr2/tldr
```

## Download

* [Distributions](https://bitbucket.org/djr2/tldr/src/master/dist/)

## Usage

```
Usage:
  -c page
      clear cache for a tldr page
      page -- Use `clearall` to clear entire cache
      -p is required if clearing cache for a specific platform
  -debug string
      enables debug logging (default "disable")
  -p platform
      platform of the tldr page
      platform -- common, linux, osx, sunos, windows (default "common")
```

### View a tldr
```
tldr <page>
```

### View a tldr for a specific platform
```
tldr -p osx <page>
```

### Clear a tldr
```
tldr -c <page>
```

### Clear a tldr for a specific platform
```
tldr -c <page> -p osx
```

### Clear entire cache
```
tldr -c clearall
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
"banner_color_1": 36,
"banner_color_2": 34,
"tldr_color": 97,
"header_color": 34,
"header_decor_color": 97,
"platform_color": 90,
"description_color": 0,
"example_color": 36,
"hypen_color": 0,
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

Pages: `https://raw.github.com/tldr-pages/tldr/master/pages/`

Zip: `https://tldr.sh/assets/tldr.zip`

To reset the configuration back to its defaults delete `config.json`
and it will be recreated. Or copy and paste the configuration from
this README above.

## License

[MIT License](https://bitbucket.org/djr2/tldr/src/master/LICENSE.md)
