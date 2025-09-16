# Viki

I just want to take my existing Obsidian Vault and turn it into a wiki I can publish statically without any fuss. All the existing Obsidian setups I saw had some extra setup involved, or made me modify my existing notes to conform to some standard I didn't care about and would have to remember if I wanted to keep using the tool, or change if I wanted to switch tools and go through it all again.

Also it seemed fun.

## Use

See `viki -h` for help, and use `-h` on any command to see available options. Viki is designed to be as simple and useful out of the box without any local modifications required.

### Generate a site

```bash
# viki generate <src> <dst>
# Example of running viki in the root directory of a vault, publishing to ./site
viki generate . ./site
```

### Serve locally

Viki can run an HTTP server locally so you can quickly see what the result would look like.

```bash
# Run in the root of the vault
viki serve
```

## Random things that need doing

### General improvements

- Add external link icon only to external links via CSS class
- Tag pages
- TOC page
    - Tag badges
- Publish to ZIP

### Page improvements

- Extract page title from main header if present
- [Code formatting](https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html)

### Sidebar improvements

- Add expand all to sidebar
- Make sidebar not shift element position when expanding

## Random ideas for hard mode

Find the latest git commit for the note and add a 'Last updated' note somewhere on the rendered page.

Look at aferox package to use github as a backing store and serve/publish straight from github repos, or even publish directly to S3.

Hot reload for serve.