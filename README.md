# Fetching GitHub Issue Data

## Done
- fetching filtered repos (13,611)
- sampling repos (500)
- fetch issues and comments (116,002)
- fixed running from nix
- fetch stars over time
- fetch commits over time
- create datasets from collected data

## Install + Run
You will have to have `nix` installed with experimental features for flakes.

You must first fetch the repos:
- `nix run .#repos` to fetch all repos that pass the filter into `./data/repos.csv`.

And sample:
- `nix run .#sample` to randomly sample 500 repos into `./data/sample.csv`

Then you can run these:
- `nix run .#comments` to fetch all the comments from sampled repos into `./data/comments.csv`.
- `nix run .#stargazers` to fetch the star history from the sampled repos into `./data/stargazers.csv`.
- `nix run .#history` to fetch the commit history from the sampled repos into `./data/history.csv`.
