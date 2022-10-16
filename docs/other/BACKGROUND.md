# Background

## Architecture

### Preserve

#### Maven

- Read conf maven.yaml.
- Get name and publicURL.
- If call to name then do the actual call to the public maven repo.
- Download the maven artifact to disk.

#### NPM

- Download the json files as .tmp.
- Replace the public URL with YAAM in .tmp files.
- Download the files via YAAM.

## Rationale for port 25213

Y is the 25th letter in the alphabet, two times 'a' equals 2 and 13 represents
the M.
