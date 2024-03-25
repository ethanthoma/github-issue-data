{ 
    pkgs ? 
        let
            inherit (builtins) fetchTree fromJSON readFile;
            inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
        in import (fetchTree nixpkgs.locked) {
            overlays = [
                (import "${fetchTree gomod2nix.locked}/overlay.nix")
            ];
        }, 
    buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
    pname = "github.com/ethanthoma/github-issue-data";
    version = "1.0.0";
    pwd = ../.;
    src = ../.;
    subPackages = [ "./cmd/comments/comments.go" ];
}
