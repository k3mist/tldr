# tldr in golang

[TLDR pages](https://tldr.sh)

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

## License

[MIT License](https://bitbucket.org/djr2/tldr/src/master/LICENSE.md)
