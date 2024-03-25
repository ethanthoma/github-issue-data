# Fetching GitHub issue data

## Done
- fetching filtered repos (13k)
- sampling repos (500) and fetching issues (8k)
- fetch comments for issues (28k)
- fixed running from nix

## TODO
- fetch stars over time
- create dataset(s) from collected data

## Install + Run
You will have to have `nix` installed with experimental features for flakes.

`nix run .#repos` to fetch all the repo ids into `./data/`.
`nix run .#comments` to fetch all the comments from 500 random repos `./data/`.
