{ name
, pkgs ? let
    inherit (builtins) fetchTree fromJSON readFile;
    inherit ((fromJSON (readFile ../flake.lock)).nodes) nixpkgs gomod2nix;
  in
  import (fetchTree nixpkgs.locked) {
    overlays = [
      (import "${fetchTree gomod2nix.locked}/overlay.nix")
    ];
  }
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
  pname = name;
  version = "1.0.0";
  src = ../.;
  pwd = ../.;
  subPackages = [ "./cmd/${name}/${name}.go" ];
}
